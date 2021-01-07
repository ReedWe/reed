// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package discover

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/reed/errors"
	"github.com/reed/log"
	"github.com/sirupsen/logrus"
	"github.com/tendermint/tmlibs/common"
	"net"
	"time"
)

var (
	packetErr = errors.New("packet message error")
)

var (
	version = []byte{0, 1}
)

const (
	pingPacket = iota + 1
	pongPacket
	findNodePacket
	findNodeRespPacket
	pingTimeout
	pongTimeout
	findNodeRespTimeout
)

const (
	responseTimeout       = 1 * time.Second
	maxFindNodeFailures   = 5
	refreshTableInterval  = 1 * time.Hour
	refreshBucketInterval = 5 * time.Minute
)

type nodeEvent uint

type UDP struct {
	common.BaseService
	conn
	listener      *udpListener
	Table         *Table
	OurNode       *Node
	timeoutEvents map[timeoutEvent]*time.Timer
	nodes         map[NodeID]*Node // record all nodes we have seen
	readCh        chan ingressPacket
	timeoutCh     chan timeoutEvent
	queryCh       chan findNodeQuery
	quitCh        chan struct{}
}

type ping struct {
	From *net.UDPAddr
	To   *net.UDPAddr
}

type pong struct {
	To *net.UDPAddr
}

type findNode struct {
	Target NodeID
}
type findNodeResp struct {
	Nodes []*rpcNode
}

type findNodeRespReply struct {
	remoteID NodeID
	nodes    []*Node
}

type ingressPacket struct {
	remoteID   NodeID
	remoteAddr *net.UDPAddr
	event      nodeEvent
	data       interface{}
	rowData    []byte
}

type timeoutEvent struct {
	node  *Node
	event nodeEvent
}

type rpcNode struct {
	IP  net.IP
	UDP uint16
	TCP uint16
	ID  NodeID
}

type conn interface {
	ReadFromUDP(b []byte) (n int, addr *net.UDPAddr, err error)
	WriteToUDP(b []byte, addr *net.UDPAddr) (n int, err error)
	Close() error
	LocalAddr() net.Addr
}

func NewDiscover() (*UDP, error) {
	// our node
	o, err := getOurNode()
	if err != nil {
		return nil, err
	}

	// kademlia table
	t, err := NewTable(o)
	if err != nil {
		return nil, err
	}

	// udp listener
	l, err := NewUDPListener(o.IP, o.UDPPort)
	if err != nil {
		return nil, err
	}

	udp := &UDP{
		OurNode:       o,
		Table:         t,
		conn:          l.conn,
		listener:      l,
		nodes:         map[NodeID]*Node{},
		timeoutEvents: map[timeoutEvent]*time.Timer{},
		readCh:        make(chan ingressPacket, 100),
		timeoutCh:     make(chan timeoutEvent),
		queryCh:       make(chan findNodeQuery),
		quitCh:        make(chan struct{}),
	}
	udp.BaseService = *common.NewBaseService(nil, "udp", udp)
	return udp, nil
}

func (u *UDP) OnStart() error {
	go u.loop()
	go u.readLoop()
	u.refresh()
	fmt.Println("★ p2p.udp Server OnStart")
	return nil
}

func (u *UDP) OnStop() {
	close(u.quitCh)
	fmt.Println("★ p2p.udp Server OnStop")
}

func (u *UDP) refresh() {
	// lookup nodes
	seeds := getSeeds()
	log.Logger.Infof("node seed count:%d", len(seeds))
	for _, seedNode := range seeds {
		if seedNode.IP.Equal(u.OurNode.IP) && seedNode.UDPPort == u.OurNode.UDPPort {
			continue
		}
		u.nodes[seedNode.ID] = seedNode
		if seedNode.state == unknown {
			transform(u, seedNode, verifyInit)
		}
		// Force-add the seed node so Lookup does something.
		// It will be deleted again if verification fails.
		u.Table.Add(seedNode)
	}
	// TODO get nodes from db?
	go u.lookup(u.OurNode.ID)
}

func (u *UDP) lookup(target NodeID) {
	log.Logger.Infof("lookup target=%s", target.ToString())
	replyCh := make(chan *findNodeRespReply, alpha)
	asked := make(map[NodeID]bool)
	pendingNodes := make(map[NodeID]*Node)
	nd := nodesByDistance{target: target}
	nd.push(u.OurNode)

	for {
		for i := 0; i < len(nd.entries) && len(pendingNodes) < alpha; i++ {
			n := nd.entries[i]
			if !asked[n.ID] {
				asked[n.ID] = true
				pendingNodes[n.ID] = n
				u.queryCh <- findNodeQuery{remote: n, target: target, reply: replyCh}
			}
		}
		// no more node
		if len(pendingNodes) == 0 {
			log.Logger.Info("no more node in pendingNodes,stop lookup")
			break
		}

		select {
		case r, ok := <-replyCh:
			log.Logger.Debugf("lookup case:replyCh")
			if ok && r != nil {
				for _, n := range r.nodes {
					if n != nil {
						log.Logger.WithFields(logrus.Fields{"remoteIP": n.IP, "port": n.UDPPort, "remoteID": n.ID.ToString(), "state": n.state}).Info("--node")
						nd.push(n)
					}
				}
				delete(pendingNodes, r.remoteID)
			} else {
				log.Logger.Debug("reply chan closed")
			}
		case <-time.After(responseTimeout):
			log.Logger.Infof("lookup timeout,pendingNodes count %d", len(pendingNodes))
			for _, v := range pendingNodes {
				if v.pendingQuery == nil {
					continue
				}
				// forget all pending requests
				close(v.pendingQuery.reply)
				// if reply is nil,don't write in replyChan
				// see func processFindNodeResp()
				v.pendingQuery.reply = nil
			}
			// start new one
			pendingNodes = make(map[NodeID]*Node)
		case <-u.quitCh:
			return
		}
	}
}

func (u *UDP) sendPong(pkt *ingressPacket) {
	log.Logger.WithFields(logrus.Fields{"remoteNodeIP": pkt.remoteAddr.IP, "port": pkt.remoteAddr.Port}).Info("send pong")
	if err := u.sendPacket(pkt.remoteAddr, pongPacket, pong{To: pkt.remoteAddr}); err != nil {
		log.Logger.Errorf("failed to pong:%v", err)
	}
}

func (u *UDP) sendPing(toNode *Node) {
	log.Logger.WithFields(logrus.Fields{"toIP": toNode.IP, "port": toNode.UDPPort}).Info("send ping")
	if err := u.sendPacket(toNode.makeUDPAddr(), pingPacket, ping{From: u.OurNode.makeUDPAddr()}); err != nil {
		log.Logger.Errorf("failed to ping:%v", err)
	}
	u.makeTimeoutEvent(toNode, pongTimeout, responseTimeout)
}

func (u *UDP) sendFindNode(toNode *Node, target NodeID) {
	log.Logger.WithFields(logrus.Fields{"toNodeIP": toNode.IP, "target": target}).Info("send findNode")
	if err := u.sendPacket(toNode.makeUDPAddr(), findNodePacket, findNode{Target: target}); err != nil {
		log.Logger.Errorf("failed to send findNode:%v", err)
	}
	u.makeTimeoutEvent(toNode, findNodeRespTimeout, responseTimeout)
}

func (u *UDP) sendFindNodeResp(toNode *Node, nd *nodesByDistance) {
	log.Logger.WithFields(logrus.Fields{"toNodeIP": toNode.IP, "Port": toNode.UDPPort, "nd": nd}).Info("->send findNodeResp")
	var rpcNodes []*rpcNode
	for _, n := range nd.entries {
		rpcNodes = append(rpcNodes, &rpcNode{
			ID:  n.ID,
			IP:  n.IP,
			UDP: n.UDPPort,
			TCP: n.TCPPort,
		})
	}
	// TODO Limit the size of the packet
	if err := u.sendPacket(toNode.makeUDPAddr(), findNodeRespPacket, findNodeResp{Nodes: rpcNodes}); err != nil {
		log.Logger.Errorf("failed to send findNodeResp:%v", err)
	}
}

func (u *UDP) sendPacket(toUDPAddr *net.UDPAddr, e nodeEvent, msg interface{}) error {
	m, err := packet(e, u.OurNode.ID, msg)
	if err != nil {
		return err
	}
	n, err := u.conn.WriteToUDP(m, toUDPAddr)
	if err != nil {
		log.Logger.Infof("failed to write to udp:%s", err.Error())
		return err
	}
	log.Logger.Tracef("nByte %d", n)
	return nil
}

func (u *UDP) loop() {
	var refreshTableTicker = time.NewTicker(refreshTableInterval)
	var refreshBucketTimer = time.NewTimer(refreshBucketInterval)
	for {
		select {
		case pkt := <-u.readCh:
			log.Logger.WithFields(logrus.Fields{"remoteIP": pkt.remoteAddr.IP, "Port": pkt.remoteAddr.Port, "event": pkt.event, "remoteID": pkt.remoteID.ToString()}).Info("->read a UDP message")
			// TODO check packet
			n := u.internNode(&pkt)
			u.processPacket(n, pkt.event, &pkt)
		case toe := <-u.timeoutCh:
			log.Logger.WithFields(logrus.Fields{"remoteID": toe.node.ID.ToString(), "IP": toe.node.IP, "PORT": toe.node.UDPPort, "event": toe.event}).Info("UDP timeout event")
			if u.timeoutEvents[toe] == nil {
				break
			}
			delete(u.timeoutEvents, toe)
			u.processPacket(toe.node, toe.event, nil)
		case f := <-u.queryCh:
			fmt.Println("<-queryCh")
			if !f.maybeExecute(u) {
				// delay execute
				f.remote.pushToDefer(&f)
			}
		case <-refreshTableTicker.C:
			// TODO if the prev refresh not done?
			log.Logger.Info("time to refresh table")
			u.refresh()
		case <-refreshBucketTimer.C:
			log.Logger.Info("time to refresh k-bucket")
			targetNode := u.Table.chooseRandomNode()
			if targetNode != nil {
				u.lookup(targetNode.ID)
			} else {
				log.Logger.Info("no target to lookup")
			}
			refreshBucketTimer.Reset(refreshBucketInterval)
		case <-u.quitCh:
			log.Logger.Info("udp.loop() quit")
			return
		default:
		}
	}
}

func (u *UDP) internNode(pkt *ingressPacket) *Node {
	if n := u.nodes[pkt.remoteID]; n != nil {
		log.Logger.Debug("remote exists in nodes")
		n.IP = pkt.remoteAddr.IP
		n.UDPPort = uint16(pkt.remoteAddr.Port)
		return n
	}
	log.Logger.Debug("remote never seen")
	n := NewNode(pkt.remoteID, pkt.remoteAddr.IP, uint16(pkt.remoteAddr.Port))
	u.nodes[pkt.remoteID] = n
	return n
}

func (u *UDP) processPong(n *Node) {
	log.Logger.WithField("nodeID", n.ID.ToString()).Info("process pong")
	u.abortTimeoutEvent(n, pongTimeout)
}

func (u *UDP) processPacket(n *Node, e nodeEvent, pkt *ingressPacket) {
	preState := n.state
	next, err := n.state.handle(u, n, e, pkt)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"nodeIP": n.IP, "UDPPort": n.UDPPort, "preStatus": preState}).Errorf("failed to handle next state:%v", err)
	}
	transform(u, n, next)
}

func (u *UDP) processQueryEvent(n *Node, e nodeEvent, pkt *ingressPacket) (*nodeState, error) {
	log.Logger.WithFields(logrus.Fields{"remoteIP": n.IP, "Port": n.UDPPort, "remoteID": n.ID.ToString(), "event": e}).Info("->process Query Event")
	switch e {
	case findNodePacket:
		// search the closest nodes with target we know
		nd := u.Table.closest(pkt.data.(*findNode).Target)
		u.sendFindNodeResp(n, nd)
		u.Table.updateConnTime(n)
		return n.state, nil
	case findNodeRespPacket:
		err := u.processFindNodeResp(n, pkt)
		return n.state, err
	case findNodeRespTimeout:
		if n.pendingQuery != nil {
			if n.pendingQuery.reply != nil {
				n.pendingQuery.reply <- nil
			}
			n.pendingQuery = nil
		}
		n.queryTimeouts++
		log.Logger.Errorf("find node response timeout,query timeout count:%d", n.queryTimeouts)
		if n.queryTimeouts >= maxFindNodeFailures && n.state == known {
			log.Logger.Error("this node timeout too much,set node state:contested")
			return contested, errors.New("node connect timeout")
		}
		return n.state, nil
	default:
		return n.state, errors.New("invalid node event")
	}
}

func (u *UDP) processFindNodeResp(n *Node, pkt *ingressPacket) error {
	log.Logger.WithFields(logrus.Fields{"remoteIP": n.IP, "Port": n.UDPPort}).Info("->processFindNodeResp")
	if n.pendingQuery == nil {
		return errors.New("no pending query")
	}
	u.abortTimeoutEvent(n, findNodeRespTimeout)

	res := pkt.data.(*findNodeResp)
	nodes := make([]*Node, len(res.Nodes))
	for i, rn := range res.Nodes {
		// check IP
		n := u.nodes[rn.ID]
		if n == nil {
			if u.OurNode.ID == rn.ID {
				log.Logger.Info("new node is self")
				continue
			}
			log.Logger.WithField("nodeID", rn.ID.ToString()).Info("new node never seen before")
			// a new node that never seen before
			n = NewNode(rn.ID, rn.IP, rn.UDP)
			if err := n.validateComplete(); err != nil {
				log.Logger.Errorf("invalid node %v", err)
				continue
			}
		} else {
			// if ip/port has change:
			if !n.IP.Equal(rn.IP) || n.UDPPort != rn.UDP || n.TCPPort != rn.TCP {
				// reject
				log.Logger.Error("node mismatch:")
				log.Logger.WithFields(logrus.Fields{"IP": n.IP, "UDPPort": n.UDPPort, "TCPPort": n.TCPPort}).Error("--local node:")
				log.Logger.WithFields(logrus.Fields{"IP": rn.IP, "UDPPort": rn.UDP, "TCPPort": rn.TCP}).Error("--remote node:")
			} else {
				log.Logger.WithField("nodeID", rn.ID.ToString()).Info("node exists")
			}
		}

		if n.state == unknown {
			log.Logger.Info("node state is unknown")
			u.nodes[n.ID] = n
			transform(u, n, verifyInit)
		}
		nodes[i] = n
	}
	// if not closed
	if n.pendingQuery.reply != nil {
		n.pendingQuery.reply <- &findNodeRespReply{remoteID: n.ID, nodes: nodes}
	}
	n.pendingQuery = nil // reset
	return nil
}

func (u *UDP) readLoop() {
	defer u.conn.Close()
	buf := make([]byte, 512)
	for {
		n, from, err := u.conn.ReadFromUDP(buf)
		if err != nil {
			log.Logger.Errorf("failed to read remoteAddr UDP:%s", err.Error())
		}
		log.Logger.WithFields(logrus.Fields{"remoteAddr": from}).Debug("read remoteAddr UDP")
		u.handlePacket(from, buf[:n])
	}
}

func (u *UDP) handlePacket(from *net.UDPAddr, buf []byte) {
	pkt := ingressPacket{
		remoteID:   BytesToHash(buf[3 : IDLength+3]),
		remoteAddr: from,
		rowData:    buf,
	}
	switch pkt.event = nodeEvent(buf[len(version) : len(version)+1][0]); pkt.event {
	case pingPacket:
		pkt.data = new(ping)
	case pongPacket:
		pkt.data = new(pong)
	case findNodePacket:
		pkt.data = new(findNode)
	case findNodeRespPacket:
		pkt.data = new(findNodeResp)
	default:
		log.Logger.Errorf("unknown packet type %d", pkt.event)
	}
	if err := json.Unmarshal(buf[23:], pkt.data); err != nil {
		log.Logger.Errorf("failed to json unmarshal %v", err)
	}
	u.readCh <- pkt
}

// makeTimeoutEvent record a timeout timer about node's event
func (u *UDP) makeTimeoutEvent(n *Node, e nodeEvent, d time.Duration) {
	te := timeoutEvent{node: n, event: e}
	u.timeoutEvents[te] = time.AfterFunc(d, func() {
		u.timeoutCh <- te
	})
}

// abortTimeoutEvent stop the timer and delete
func (u *UDP) abortTimeoutEvent(n *Node, e nodeEvent) {
	te := u.timeoutEvents[timeoutEvent{node: n, event: e}]
	if te != nil {
		te.Stop()
		delete(u.timeoutEvents, timeoutEvent{node: n, event: e})
	}
}

func (ns *nodeState) canQuery() bool {
	return ns == known || ns == contested || ns == unresponsive
}

// header
// [2]byte	version
// [1]byte	packetType
// [20]byte	nodeId
func packet(e nodeEvent, ourId NodeID, msg interface{}) ([]byte, error) {
	b := new(bytes.Buffer)
	b.Write(version)
	b.WriteByte(byte(e))
	b.Write(ourId[:])
	ms, err := json.Marshal(msg)
	if err != nil {
		log.Logger.Infof("failed to packet msg:%s", err.Error())
		return nil, errors.Wrapf(packetErr, err.Error())
	}
	b.Write(ms)
	return b.Bytes(), nil
}

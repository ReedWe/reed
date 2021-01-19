// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package p2p

import (
	"bufio"
	"github.com/reed/errors"
	"github.com/reed/log"
	"github.com/reed/p2p/discover"
	"github.com/sirupsen/logrus"
	"github.com/tendermint/tmlibs/common"
	"net"
	"sync"
	"time"
)

const (
	dialTimeout      = 3 * time.Second
	handshakeTimeout = 5 * time.Second

	peerUpdateInterval = time.Second * 10

	connectionSize = 30
)

var (
	addPeerFromAcceptErr = errors.New("failed to add peer from accept")
	dialErr              = errors.New("failed to dial")
)

type Network struct {
	common.BaseService
	pm          *PeerMap
	dialing     *PeerDialing
	table       *discover.Table
	ourNodeInfo *NodeInfo
	handlerServ Handler
	acceptCh    <-chan net.Conn
	disConnCh   chan string
	quitCh      chan struct{}
}

func NewNetWork(ourNode *discover.Node, t *discover.Table, acceptCh <-chan net.Conn, handlerServ Handler) (*Network, error) {
	n := &Network{
		pm:          NewPeerMap(),
		dialing:     NewPeerDialing(),
		table:       t,
		ourNodeInfo: NewOurNodeInfo(ourNode.ID, ourNode.IP, ourNode.TCPPort),
		handlerServ: handlerServ,
		acceptCh:    acceptCh,
		disConnCh:   make(chan string),
		quitCh:      make(chan struct{}),
	}
	n.BaseService = *common.NewBaseService(nil, "network", n)
	return n, nil
}

func (n *Network) OnStart() error {
	go n.loop()
	go n.loopFillPeer()
	log.Logger.Info("★★p2p.TCP Network Server OnStart")
	return nil
}

func (n *Network) OnStop() {
	close(n.quitCh)
	for _, v := range n.pm.peers {
		n.releasePeer(v)
	}
	close(n.disConnCh)
	log.Logger.Info("★★p2p.TCP Network Server OnStop")
}

func (n *Network) loop() {
	for {
		select {
		case c, ok := <-n.acceptCh:
			if !ok {
				return
			}
			log.Logger.WithFields(logrus.Fields{"remote addr": c.RemoteAddr().String()}).Info("->accept peer")
			if err := n.addPeerFromAccept(c); err != nil {
				log.Logger.Error(err)
			}
		case addr := <-n.disConnCh:
			log.Logger.WithField("peerAddr", addr).Info("remove disconnection peer")
			peer := n.pm.get(addr)
			if peer != nil {
				n.releasePeer(peer)
			}
		case <-n.quitCh:
			log.Logger.Info("[loop] quit")
			return
		}
	}
}

func (n *Network) loopFillPeer() {
	timer := time.NewTimer(peerUpdateInterval)
	for {
		select {
		case <-timer.C:
			n.fillPeer()
			timer.Reset(peerUpdateInterval)
		case <-n.quitCh:
			log.Logger.Info("[loopFillPeer] quit")
			return
		}
	}
}

func (n *Network) addPeerFromAccept(conn net.Conn) error {
	if n.pm.peerCount() >= connectionSize {
		_ = conn.Close()
		return errors.Wrap(addPeerFromAcceptErr, "enough peers")
	}
	return n.connectPeer(conn)
}

func (n *Network) fillPeer() {
	log.Logger.Debug("time to fill Peer")
	nodes := n.table.GetWithExclude(activelyPeerCount, n.pm.IDs())
	if len(nodes) == 0 {
		log.Logger.Debug("no available node to dial")
		return
	}

	var wg sync.WaitGroup
	var waitForProcess []string
	for _, node := range nodes {
		addr := toAddress(node.IP, node.TCPPort)
		if n.dialing.exist(addr) {
			log.Logger.WithField("peerAddr", addr).Info("peer is dialing")
			continue
		} else {
			waitForProcess = append(waitForProcess, addr)
		}
	}

	wg.Add(len(waitForProcess))
	for _, addr := range waitForProcess {
		go n.dialPeerAndAdd(addr, &wg)
	}
	wg.Wait()
}

func (n *Network) dialPeerAndAdd(addr string, wg *sync.WaitGroup) {
	defer func() {
		n.dialing.remove(addr)
		wg.Done()
	}()

	n.dialing.add(addr)
	rawConn, err := n.dial(addr)
	if err != nil {
		log.Logger.WithField("peerAddr", addr).Errorf("failt to dial peer:%v", err)
		return
	}
	if err = n.connectPeer(rawConn); err != nil {
		log.Logger.WithField("peerAddr", addr).Errorf("failed to connect peer: %v", err)
		return
	}
}

func (n *Network) connectPeer(rawConn net.Conn) error {
	nodeInfo, err := n.handshake(rawConn)
	if err != nil {
		return err
	}

	peer := NewPeer(n.ourNodeInfo, nodeInfo, n.disConnCh, rawConn, n.handlerServ)
	if err = peer.Start(); err != nil {
		return err
	}
	n.pm.add(peer)
	log.Logger.WithField("peerAddr", peer.nodeInfo.RemoteAddr).Info("Add a new peer and started.")
	return nil
}

func (n *Network) releasePeer(peer *Peer) {
	if err := peer.Stop(); err != nil {
		log.Logger.WithField("peerAddr", peer.nodeInfo.RemoteAddr).Error(err)
	}
	n.pm.remove(peer.nodeInfo.RemoteAddr)
}

func (n *Network) dial(address string) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", address, dialTimeout)
	if err != nil {
		return nil, errors.Wrap(dialErr, err)
	}
	return conn, nil
}

func (n *Network) handshake(conn net.Conn) (*NodeInfo, error) {
	log.Logger.WithFields(logrus.Fields{"self": n.ourNodeInfo.RemoteAddr, "remote": conn.RemoteAddr()}).Debug("ready to handshake")
	if err := conn.SetDeadline(time.Now().Add(handshakeTimeout)); err != nil {
		return nil, err
	}
	if err := write(conn, []byte{handshakeCode}); err != nil {
		return nil, err
	}

	input := bufio.NewScanner(conn)
	for input.Scan() {
		switch b := input.Bytes(); b[0] {
		case handshakeCode:
			log.Logger.WithFields(logrus.Fields{"self": n.ourNodeInfo.RemoteAddr, "remote": conn.RemoteAddr()}).Debug("process [handshakeCode]")
			if err := writeOurNodeInfo(conn, n.ourNodeInfo); err != nil {
				log.Logger.Error(err)
			}
		case handshakeRespCode:
			log.Logger.WithFields(logrus.Fields{"self": n.ourNodeInfo.RemoteAddr, "remote": conn.RemoteAddr()}).Debug("process [handshakeRespCode]")
			return NewNodeInfoFromBytes(b[1:], conn)
		}
	}
	if input.Err() != nil {
		return nil, input.Err()
	}
	return nil, errors.New("handshake error")
}

// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package discover

import (
	"encoding/hex"
	"github.com/reed/blockchain/config"
	"github.com/reed/errors"
	"github.com/reed/log"
	"net"
)

var (
	getOurNodeErr = errors.New("build our node error")
)

const (
	IDLength = 32
)

type NodeID [IDLength]byte

type Node struct {
	ID      NodeID
	IP      net.IP
	TCPPort uint16
	UDPPort uint16
	nodeNetGuts
}

type nodeNetGuts struct {
	state         *nodeState
	deferQueries  []*findNodeQuery
	pendingQuery  *findNodeQuery
	queryTimeouts int
}

type findNodeQuery struct {
	remote *Node
	target NodeID
	reply  chan<- *findNodeRespReply
}

func (n *Node) makeUDPAddr() *net.UDPAddr {
	return &net.UDPAddr{IP: n.IP, Port: int(n.UDPPort)}
}

func (n *Node) validateComplete() error {
	if n.IP == nil {
		return errors.New("node IP nil")
	}
	if n.UDPPort == 0 {
		return errors.New("missing UDP Port")
	}
	if n.TCPPort == 0 {
		return errors.New("missing TCP Port")
	}
	if n.IP.IsMulticast() {
		return errors.New("invalid IP multicast")
	}
	if n.IP.IsUnspecified() {
		return errors.New("invalid IP unspecified")
	}
	return nil
}

func NewNode(id NodeID, ip net.IP, port uint16) *Node {
	n := &Node{
		ID:      id,
		IP:      ip,
		UDPPort: port,
		TCPPort: port,
		nodeNetGuts: nodeNetGuts{
			state: unknown,
		},
	}

	return n
}

func getOurNode() (*Node, error) {
	node := resolveNode(config.Default.OurNode)
	if node == nil {
		return nil, errors.Wrapf(getOurNodeErr, "config.toml:missing self node")
	}
	log.Logger.Info("Our Node Info:")
	log.Logger.Infof("--ID:%s", node.ID.ToString())
	log.Logger.Infof("--IP:%s", node.IP.String())
	log.Logger.Infof("--UDP Port:%d", node.UDPPort)
	log.Logger.Infof("--TCP Port:%d", node.TCPPort)
	return node, nil
}

func BytesToHash(b []byte) NodeID {
	var h NodeID
	h.SetBytes(b)
	return h
}

func (h *NodeID) ToString() string {
	return hex.EncodeToString(h.Bytes())
}

func (h NodeID) Bytes() []byte { return h[:] }

func (h *NodeID) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-IDLength:]
	}

	copy(h[IDLength-len(b):], b)
}

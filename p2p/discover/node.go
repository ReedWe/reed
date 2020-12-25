// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package discover

import (
	"encoding/hex"
	"github.com/reed/blockchain/config"
	"github.com/reed/errors"
	"net"
	"strconv"
	"strings"
)

var (
	getOurNodeErr = errors.New("build our node error")
)

const (
	IDLength = 20
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

func NewNode(id NodeID, ip net.IP, udpPort uint16) *Node {
	n := &Node{
		ID:      id,
		IP:      ip,
		UDPPort: udpPort,
		nodeNetGuts: nodeNetGuts{
			state: unknown,
		},
	}

	return n
}

func getOurNode() (*Node, error) {
	arr := strings.Split(config.Default.LocalAddr, ":")
	addr, err := net.ResolveIPAddr("ip", arr[0])
	if err != nil {
		return nil, errors.Wrapf(getOurNodeErr, "failed to resolve our node IP address")
	}

	port, err := strconv.ParseUint(arr[1], 10, 16)
	if err != nil {
		return nil, errors.Wrapf(getOurNodeErr, "failed to resolve our node IP port")
	}

	//TODO our node ID
	ds, _ := hex.DecodeString(config.Default.OurID)
	return NewNode(BytesToHash(ds), addr.IP, uint16(port)), nil
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

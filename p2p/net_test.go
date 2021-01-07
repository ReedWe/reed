// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package p2p

import (
	"fmt"
	"github.com/reed/blockchain/netsync"
	"github.com/reed/log"
	"github.com/reed/p2p/discover"
	"net"
	"testing"
)

var (
	port      = 8000
	ourNodeID = discover.NodeID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(5), byte(4), byte(3), byte(2), byte(1)}
	ourAddr   = "127.0.0.1:8000"

	// tos = the other side
	tosPort   = 7000
	tosNodeID = discover.NodeID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(255), byte(255), byte(255), byte(255), byte(255), byte(255)}
	tosAddr   = "127.0.0.1:7000"
)

func getOurNode() *discover.Node {
	return &discover.Node{
		ID:      ourNodeID,
		TCPPort: uint16(port),
		UDPPort: uint16(port),
		IP:      net.IP{127, 0, 0, 1},
	}
}

func getRemoteNode() *discover.Node {
	return &discover.Node{
		ID:      tosNodeID,
		TCPPort: uint16(tosPort),
		UDPPort: uint16(tosPort),
		IP:      net.IP{127, 0, 0, 1},
	}
}

func getOurNodeInfo() *NodeInfo {
	return &NodeInfo{
		ID:         ourNodeID,
		RemoteAddr: ourAddr,
	}
}

func startUpTheOtherSidePeer() error {
	listen, err := net.Listen("tcp", tosAddr)
	if err != nil {
		return err
	}

	other := mockKadTableRemote()
	netWork, err := NewNetWork(getRemoteNode(), other, nil, netsync.Handle)

	go func() {
		for {
			conn, err2 := listen.Accept()
			if err2 != nil {
				fmt.Println(err2)
			}
			// handshake
			netWork.addPeerFromAccept(conn)
		}
	}()
	return nil
}

func TestFillPeer(t *testing.T) {
	log.Init()
	table := mockKadTableOur()
	netWork, err := NewNetWork(getOurNode(), table, nil, netsync.Handle)

	startUpTheOtherSidePeer()

	if err != nil {
		t.Fatal(err)
	}
	netWork.fillPeer()
	if netWork.pm.peerCount() != 1 {
		t.Fatal("peer count error")
	}
	p := netWork.pm.get(tosAddr)
	if p.nodeInfo.RemoteAddr != tosAddr {
		t.Fatal("tos peer in map not right")
	}
}

func mockKadTableOur() *discover.Table {
	self := &discover.Node{
		ID: ourNodeID,
	}
	tb, _ := discover.NewTable(self)

	n := &discover.Node{
		ID:      tosNodeID,
		TCPPort: uint16(tosPort),
		UDPPort: uint16(tosPort),
		IP:      net.IP{127, 0, 0, 1},
	}

	tb.Add(n)
	return tb
}

func mockKadTableRemote() *discover.Table {
	self := &discover.Node{
		ID: tosNodeID,
	}
	tb, _ := discover.NewTable(self)

	n := &discover.Node{
		ID:      ourNodeID,
		TCPPort: uint16(port),
		UDPPort: uint16(port),
		IP:      net.IP{127, 0, 0, 1},
	}

	tb.Add(n)
	return tb
}

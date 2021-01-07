// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package discover

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net"
	"testing"
)

func TestContains(t *testing.T) {
	id1, _ := hex.DecodeString("7b52009b64fd0a2a49e6d8a939753077792b0554")
	id2, _ := hex.DecodeString("40bd001563085fc35165329ea1ff5c5ecbdbbeef")
	id3, _ := hex.DecodeString("7110eda4d09e062aa5e4a390b0a572ac0d2c0220")

	var ns []*bNode
	ns = append(ns, &bNode{node: &Node{ID: BytesToHash(id1)}})
	ns = append(ns, &bNode{node: &Node{ID: BytesToHash(id2)}})

	if !contains(ns, BytesToHash(id1)) {
		t.Error("contains error (expect id in ns)")
	}

	if contains(ns, BytesToHash(id3)) {
		t.Error("contains error (expect ns does not exist request id)")
	}
}

func TestComputeDist(t *testing.T) {
	ta := NodeID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(3)}
	id1 := NodeID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(3)}
	id2 := NodeID{0, 0, 0, 0, 0, byte(8), 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(2)}

	dist := computeDist(ta, id1, id2)
	fmt.Println(dist)
}

func TestNodesByDistance(t *testing.T) {
	nd := nodesByDistance{
		target: NodeID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(3)},
		entries: []*Node{
			{ID: NodeID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(1)}},
			{ID: NodeID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(2), 0}},
			{ID: NodeID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(12), 0, 0, 0, 0, 0, 0, byte(2), 0, 0}},
		},
	}
	node := &Node{
		ID: NodeID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(2), 0, byte(1)},
	}
	nd.push(node)

	if len(nd.entries) != 4 {
		t.Error("failed to push node")
	}

	if !bytes.Equal(nd.entries[2].ID.Bytes(), node.ID.Bytes()) {
		t.Error("nodesByDistance push error")
	}
}

type tn struct {
	name string
	node *Node
}

func TestGetWithExclude(t *testing.T) {
	tb := newTable()

	minDist := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(2), byte(1)}
	secondDist := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(3), byte(2), byte(1)}

	ns := tb.GetWithExclude(1, nil)
	if !bytes.Equal(ns[0].ID.Bytes(), minDist) {
		t.Fatal("the first(minimum distance) not right")
	}

	ns2 := tb.GetWithExclude(4, []string{net.IP{123, 123, 123, 13}.String() + ":" + "8002"})
	if len(ns2) != 4 {
		t.Fatal("wrong count")
	}
	for _, n := range ns2 {
		if bytes.Equal(n.ID.Bytes(), secondDist) {
			t.Fatal("does not exclude the given node")
		}
	}

}

func newTable() *Table {
	our := &Node{
		ID: NodeID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(1)},
	}
	tb, _ := NewTable(our)

	n2 := &Node{
		ID: NodeID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(2), byte(1)},
	}
	n3 := &Node{
		ID:      NodeID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(3), byte(2), byte(1)},
		TCPPort: 8002,
		UDPPort: 8001,
		IP:      net.IP{123, 123, 123, 13},
	}
	n4 := &Node{
		ID: NodeID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(4), byte(3), byte(2), byte(1)},
	}
	n5 := &Node{
		ID: NodeID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(5), byte(4), byte(3), byte(2), byte(1)},
	}
	n6 := &Node{
		ID: NodeID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(6), byte(5), byte(4), byte(3), byte(2), byte(1)},
	}
	n7 := &Node{
		ID: NodeID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(7), byte(6), byte(5), byte(4), byte(3), byte(2), byte(1)},
	}
	var tns []*tn
	tns = append(tns, &tn{
		name: "n2",
		node: n2,
	})
	tns = append(tns, &tn{
		name: "n3",
		node: n3,
	})
	tns = append(tns, &tn{
		name: "n4",
		node: n4,
	})
	tns = append(tns, &tn{
		name: "n5",
		node: n5,
	})
	tns = append(tns, &tn{
		name: "n6",
		node: n6,
	})
	tns = append(tns, &tn{
		name: "n7",
		node: n7,
	})

	for _, t := range tns {
		tb.Add(t.node)
	}

	for i, b := range tb.Bucket {
		for _, n := range b {
			for _, t := range tns {
				if n.node == t.node {
					fmt.Printf("name:%s k-bucket:%d\n", t.name, i)
				}
			}
		}
	}
	return tb
}

// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package discover

import (
	"bytes"
	"encoding/hex"
	"fmt"
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

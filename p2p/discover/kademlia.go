// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package discover

import (
	"github.com/reed/log"
	"github.com/reed/types"
	"math/rand"
	"sort"
	"sync"
	"time"
)

// References
// Kademlia: A Peer-to-peer Information System Based on the XOR Metric
// http://www.scs.stanford.edu/~dm/home/papers/kpos.pdf

const (
	kBucketSize = 16

	// The most important procedure a Kademlia participant must perform is to
	// locate the k closest nodes to some given node ID. We call this procedure
	// a node lookup. Kademlia employs a recursive algorithm for node lookups.
	// The lookup initiator starts by picking α nodes remoteAddr its closest non-empty
	// k-bucket (or, if that bucket has fewer than α entries, it just takes the
	// α closest nodes it knows of). The initiator then sends parallel, asynchronous
	// FIND NODE RPCs to the α nodes it has chosen.
	alpha = 3

	IDBits = len(types.Hash{}) * 8
)

type Table struct {
	mutex   sync.Mutex
	Bucket  [IDBits][]*bNode
	OurNode *Node
}

type bNode struct {
	node       *Node
	lastConnAt time.Time
}

func NewTable(ourNode *Node) (*Table, error) {
	t := &Table{
		Bucket:  [IDBits][]*bNode{},
		OurNode: ourNode,
	}
	return t, nil
}

func (t *Table) getNodeAccurate(id NodeID) *bNode {
	kbs := t.Bucket[logarithmDist(t.OurNode.ID, id)]
	bn := getNodeFromKbs(kbs, id)
	if bn == nil {
		return nil
	}
	return bn
}

func (t *Table) delete(id NodeID) {
	dist := logarithmDist(t.OurNode.ID, id)
	for i, bn := range t.Bucket[dist] {
		if bn.node.ID == id {
			t.Bucket[dist] = append(t.Bucket[dist][:i], t.Bucket[dist][i+1:]...)
			return
		}
	}
	log.Logger.Info("delete node complete")
	t.printLog()
}

func (t *Table) Add(n *Node) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if n == t.OurNode {
		return
	}
	for _, b := range t.Bucket {
		if contains(b, n.ID) {
			log.Logger.Debug("node exists in table")
			return
		}
	}

	dist := logarithmDist(t.OurNode.ID, n.ID)

	if len(t.Bucket[dist]) < kBucketSize {
		t.Bucket[dist] = append(t.Bucket[dist], &bNode{node: n, lastConnAt: time.Now().UTC()})
	}

	log.Logger.Info("add node complete")
	t.printLog()

	// TODO when len(kBucket) >= kBucketSize
	// do something...
}

func (t *Table) updateConnTime(n *Node) {
	bn := t.getNodeAccurate(n.ID)
	if bn != nil {
		bn.lastConnAt = time.Now().UTC()
	}
}

func (t *Table) closest(target NodeID) *nodesByDistance {
	nd := &nodesByDistance{target: target}
	for _, b := range t.Bucket {
		for _, n := range b {
			nd.push(n.node)
		}
	}
	return nd
}

// chooseRandomNode choose the node who has not performed a node lookup within an hour.
func (t *Table) chooseRandomNode() *Node {
	bt := time.Now().Add(-1 * time.Hour)

	n := t.getRandomOne(func(bn *bNode) bool {
		return bt.After(bn.lastConnAt)
	})
	if n == nil {
		n = t.getRandomOne(func(bn *bNode) bool {
			return true
		})
	}
	return n
}

func (t *Table) getRandomOne(condition func(bn *bNode) bool) *Node {
	var nodes []*Node
	for _, b := range t.Bucket {
		for _, n := range b {
			if condition(n) {
				nodes = append(nodes, n.node)
			}
		}
		if len(nodes) >= kBucketSize {
			break
		}
	}
	if len(nodes) == 0 {
		return nil
	}
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return nodes[rd.Intn(len(nodes))]
}

func (t *Table) GetWithExclude(count int, excludePeerIDs []NodeID) []*Node {
	exclude := func(id NodeID) bool {
		if len(excludePeerIDs) == 0 {
			return false
		}
		for _, eID := range excludePeerIDs {
			if eID == id {
				return true
			}
		}
		return false
	}

	var nodes []*Node

loop:
	for _, b := range t.Bucket {
		for _, n := range b {
			if !exclude(n.node.ID) {
				nodes = append(nodes, n.node)
				if count == len(nodes) {
					break loop
				}
			}
		}
	}
	return nodes
}

func (t *Table) printLog() {
	for i, b := range t.Bucket {
		if len(b) == 0 {
			continue
		}
		log.Logger.Debugf("K-Bucket No:%d", i)
		for _, bn := range b {
			log.Logger.Debugf("---Addr:%s:%d ID:%s LastConnAt:%v", bn.node.IP, bn.node.UDPPort, bn.node.ID.ToString(), bn.lastConnAt)
		}
	}
}

type nodesByDistance struct {
	entries []*Node
	target  NodeID
}

func (nd *nodesByDistance) push(n *Node) {
	for _, entry := range nd.entries {
		if entry.ID == n.ID {
			return
		}
	}
	ix := sort.Search(len(nd.entries), func(i int) bool {
		return computeDist(nd.target, nd.entries[i].ID, n.ID) > 0
	})
	if len(nd.entries) < kBucketSize {
		nd.entries = append(nd.entries, n)
	}
	if ix == len(nd.entries) {
		// farther away than all nodes we already have.
		// if there was room for it, the node is now the last element.
	} else {
		// slide existing entries down to make room
		// this will overwrite the entry we just appended.
		copy(nd.entries[ix+1:], nd.entries[ix:])
		nd.entries[ix] = n
	}
}

func computeDist(target, a, b NodeID) int {
	for i := range target {
		da := a[i] ^ target[i]
		db := b[i] ^ target[i]
		if da > db {
			return 1
		} else if da < db {
			return -1
		}
	}
	return 0
}

func getNodeFromKbs(ns []*bNode, id NodeID) *bNode {
	if len(ns) == 0 {
		return nil
	}
	for _, bn := range ns {
		if bn.node.ID == id {
			return bn
		}
	}
	return nil
}

func contains(ns []*bNode, id NodeID) bool {
	return getNodeFromKbs(ns, id) != nil
}

// logarithmDist return distance between a and b
// return log2(a^b)

//	k-bucket	distance	description
//	0			[2^0,2^1)	存放距离为1，且前255bit相同，第256bit开始不同（即前255bit为0）
//	1			[2^1,2^2)	存放距离为2~3，且前254bit相同，第255bit开始不同
//	2			[2^2,2^3)	存放距离为4~7，且前253bit相同，第254bit开始不同
//	...
//	MEMO:
//	ID长度为32Byte，256bit。
//	上面循环每一位，进行异或（^）操作，结果0表示相同，1表示不同
//	所以“前导0个数为255”表示有255个bit是相同的
func logarithmDist(a, b NodeID) int {
	for i := range a {
		x := a[i] ^ b[i]
		if x != 0 {
			lz := i*8 + lzcount[x] // 256bit leading zero counts
			return IDBits - 1 - lz
		}
	}
	return 0
}

// table of leading zero counts for bytes [0..255]
var lzcount = [256]int{
	8, 7, 6, 6, 5, 5, 5, 5,
	4, 4, 4, 4, 4, 4, 4, 4,
	3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3,
	2, 2, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 2, 2,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

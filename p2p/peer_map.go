// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package p2p

import (
	"github.com/reed/p2p/discover"
	"sync"
)

const (
	activelyPeerCount = 15
)

type PeerMap struct {
	mtx   sync.RWMutex
	peers map[string]*Peer
}

func NewPeerMap() *PeerMap {
	return &PeerMap{
		peers: map[string]*Peer{},
	}
}

func (pl *PeerMap) add(peer *Peer) {
	defer pl.mtx.Unlock()
	pl.mtx.Lock()
	pl.peers[peer.nodeInfo.RemoteAddr] = peer
}

func (pl *PeerMap) remove(addr string) {
	defer pl.mtx.Unlock()
	pl.mtx.Lock()

	delete(pl.peers, addr)
}

func (pl *PeerMap) get(addr string) *Peer {
	defer pl.mtx.RUnlock()
	pl.mtx.RLock()
	return pl.peers[addr]
}

func (pl *PeerMap) peerCount() int {
	defer pl.mtx.RUnlock()
	pl.mtx.RLock()
	return len(pl.peers)
}

func (pl *PeerMap) IDs() []discover.NodeID {
	defer pl.mtx.RUnlock()
	pl.mtx.RLock()
	if len(pl.peers) == 0 {
		return nil
	}

	var ids []discover.NodeID
	for _, p := range pl.peers {
		ids = append(ids, p.nodeInfo.ID)
	}
	return ids
}

// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package p2p

import "sync"

type PeerDialing struct {
	mtx   sync.RWMutex
	peers map[string]struct{}
}

func NewPeerDialing() *PeerDialing {
	return &PeerDialing{
		peers: map[string]struct{}{},
	}
}

func (pd *PeerDialing) exist(addr string) bool {
	defer pd.mtx.RUnlock()
	pd.mtx.RLock()
	_, ok := pd.peers[addr]
	return ok
}

func (pd *PeerDialing) add(addr string) {
	defer pd.mtx.Unlock()
	pd.mtx.Lock()
	pd.peers[addr] = struct{}{}
}

func (pd *PeerDialing) remove(addr string) {
	defer pd.mtx.Unlock()
	pd.mtx.Lock()
	delete(pd.peers, addr)
}

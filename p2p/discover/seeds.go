// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package discover

import (
	"encoding/hex"
	"net"
)

func getSeeds() []*Node {
	return []*Node{
		getSeed(),
	}
}

func getSeed() *Node {
	sid, _ := hex.DecodeString("67032b2b262d837fbe2a0608409986c571350689")
	addr, _ := net.ResolveIPAddr("ip", "127.0.0.1")
	return NewNode(BytesToHash(sid), addr.IP, 30398)
}

// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package p2p

import "net"

type Node struct {
	ID      [32]byte
	IP      net.IP
	TCPPort uint16
	UDPPort uint16
}

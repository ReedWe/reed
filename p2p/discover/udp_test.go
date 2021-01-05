// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package discover

import (
	"fmt"
	"net"
	"testing"
)

func TestPacket(t *testing.T) {
	m := make(map[*timeoutEvent]string)

	addr, _ := net.ResolveIPAddr("ip", "127.0.0.1")

	n1 := &Node{
		ID:      NodeID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(3)},
		IP:      addr.IP,
		UDPPort: 9888,
		nodeNetGuts: nodeNetGuts{
			state:        unknown,
			deferQueries: []*findNodeQuery{},
		},
	}

	t1 := timeoutEvent{
		node:  n1,
		event: findNodeRespPacket,
	}

	t2 := timeoutEvent{
		node: NewNode(
			NodeID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(2), 0, 0, 0, 0, byte(3)},
			addr.IP,
			9889,
		),
		event: findNodeRespPacket,
	}

	m[&t1] = "1111111111111111111111"
	m[&t2] = "2222222222222222222222"

	// nt := timeoutEvent{
	//	node: NewNode(
	//		NodeID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(3)},
	//		addr.IP,
	//		9888,
	//	),
	//	event: findNodeRespPacket,
	// }

	// n := m[timeoutEvent{
	//	node:  n1,
	//	event: findNodeRespPacket,
	// }]

	fmt.Println(findNodeRespTimeout)

}

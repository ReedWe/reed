// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package p2p

import (
	"net"
	"strconv"
)

type udpListener struct {
	conn *net.UDPConn
}

func NewUDPListener(address string, port int) (*udpListener, error) {
	listener := &udpListener{}
	addr, err := net.ResolveUDPAddr("udp", address+":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	c, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	listener.conn = c
	return listener, nil
}

// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package discover

import (
	"net"
)

type udpListener struct {
	conn *net.UDPConn
}

func NewUDPListener(ip net.IP, port uint16) (*udpListener, error) {
	listener := &udpListener{}
	c, err := net.ListenUDP("udp", &net.UDPAddr{IP: ip, Port: int(port)})
	if err != nil {
		return nil, err
	}

	listener.conn = c
	return listener, nil
}

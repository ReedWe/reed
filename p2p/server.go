// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package p2p

import "github.com/reed/p2p/discover"

type Server struct {
	tcpListener *Listener
	udp         *discover.UDP
}

func NewP2PServer() (*Server, error) {
	udp, err := discover.NewDiscover()
	if err != nil {
		return nil, err
	}

	listener, err := NewListener(udp.OurNode.IP, udp.OurNode.TCPPort)
	if err != nil {
		return nil, err
	}

	return &Server{
		tcpListener: listener,
		udp:         udp,
	}, nil
}

func (s *Server) Start() {
	s.udp.Start()
	s.tcpListener.Start()
}

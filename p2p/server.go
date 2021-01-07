// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package p2p

import (
	"github.com/reed/blockchain/netsync"
	"github.com/reed/p2p/discover"
	"github.com/tendermint/tmlibs/common"
)

type Server struct {
	common.BaseService
	tcpListener *Listener
	udp         *discover.UDP
	network     *Network
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

	network, err := NewNetWork(udp.OurNode, udp.Table, listener.acceptCh, netsync.Handle)
	if err != nil {
		return nil, err
	}

	serv := &Server{
		tcpListener: listener,
		udp:         udp,
		network:     network,
	}
	serv.BaseService = *common.NewBaseService(nil, "p2pserver", serv)
	return serv, nil
}

func (s *Server) OnStart() error {
	if err := s.udp.Start(); err != nil {
		return err
	}
	if err := s.tcpListener.Start(); err != nil {
		return err
	}
	if err := s.network.Start(); err != nil {
		return err
	}
	return nil
}

func (s *Server) OnStop() {
	s.udp.Stop()
	s.tcpListener.Stop()
	s.network.Stop()
}

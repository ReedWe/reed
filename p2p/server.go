// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package p2p

import (
	"github.com/reed/log"
	"github.com/reed/p2p/discover"
	"github.com/tendermint/tmlibs/common"
)

type Server struct {
	common.BaseService
	tcpListener *Listener
	udp         *discover.UDP
	network     *Network
}

func NewP2PServer(handleService Handler) (*Server, error) {
	udp, err := discover.NewDiscover()
	if err != nil {
		return nil, err
	}

	listener, err := NewListener(udp.OurNode.IP, udp.OurNode.TCPPort)
	if err != nil {
		return nil, err
	}

	network, err := NewNetWork(udp.OurNode, udp.Table, listener.acceptCh, handleService)
	if err != nil {
		return nil, err
	}

	serv := &Server{
		tcpListener: listener,
		udp:         udp,
		network:     network,
	}
	serv.BaseService = *common.NewBaseService(nil, "p2pServer", serv)
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
	if err := s.udp.Stop(); err != nil {
		log.Logger.Errorf("failed to stop UDP Server:%v", err)
	}
	if err := s.tcpListener.Stop(); err != nil {
		log.Logger.Errorf("failed to stop Listener Server:%v", err)
	}
	if err := s.network.Stop(); err != nil {
		log.Logger.Errorf("failed to stop TCP Newwork Server:%v", err)
	}
}

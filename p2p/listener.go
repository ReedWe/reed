// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package p2p

import (
	"github.com/reed/log"
	"github.com/tendermint/tmlibs/common"
	"net"
	"strconv"
)

const (
	acceptBufSize = 10
)

type Listener struct {
	common.BaseService
	listen   net.Listener
	acceptCh chan net.Conn
}

func NewListener(ip net.IP, port uint16) (*Listener, error) {
	l, err := net.Listen("tcp", net.JoinHostPort(ip.String(), strconv.FormatUint(uint64(port), 10)))
	if err != nil {
		return nil, err
	}
	listener := &Listener{
		listen:   l,
		acceptCh: make(chan net.Conn, acceptBufSize),
	}
	listener.BaseService = *common.NewBaseService(nil, "listener", listener)
	return listener, nil
}

func (l *Listener) OnStart() error {
	go l.loop()
	log.Logger.Info("★★p2p.Listener Server OnStart")
	return nil
}

func (l *Listener) OnStop() {
	if err := l.listen.Close(); err != nil {
		log.Logger.Errorf("failed to stop listener:%v", err)
	}
	log.Logger.Info("★★p2p.Listener Server OnStop")
}

func (l *Listener) loop() {
	for {
		c, err := l.listen.Accept()
		if err != nil {
			log.Logger.Errorf("listener.loop %v", err)
			break
		}
		l.acceptCh <- c
	}
	close(l.acceptCh)
}

// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package p2p

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/reed/log"
	"github.com/tendermint/tmlibs/common"
	"net"
	"sync"
)

const (
	connectSt = iota + 1
	localDisConnSt
	remoteDisConnSt
)

type connState uint

type HandleFunc func(msg []byte) []byte

type Conn struct {
	common.BaseService
	ourNodeInfo *NodeInfo
	peerAddr    string
	rawConn     net.Conn
	state       state
	disConnCh   chan<- string
	handle      HandleFunc
}

type state struct {
	mtx sync.RWMutex
	val connState
}

func NewConnection(peerAddr string, disConnCh chan<- string, rawConn net.Conn, ourNodeInfo *NodeInfo, handleFunc HandleFunc) *Conn {
	conn := &Conn{
		ourNodeInfo: ourNodeInfo,
		peerAddr:    peerAddr,
		rawConn:     rawConn,
		disConnCh:   disConnCh,
		state:       state{val: connectSt},
		handle:      handleFunc,
	}
	conn.BaseService = *common.NewBaseService(nil, "conn", conn)
	return conn
}

func (c *Conn) OnStart() error {
	go c.readGoroutine()
	return nil
}

func (c *Conn) OnStop() {
	// flag local set disconnection,don't need send disConnCh.
	// see func readGoroutine().
	c.setState(localDisConnSt)
	if err := c.rawConn.Close(); err != nil {
		log.Logger.Error(err)
	}
}

func (c *Conn) readGoroutine() {
	input := bufio.NewScanner(c.rawConn)
	for input.Scan() {
		if c.specialMsg(input.Bytes()) {
			continue
		}
		writeMsg := c.handle(input.Bytes())
		if writeMsg != nil {
			if err := c.write(writeMsg); err != nil {
				log.Logger.Error("failed to write:%v", err)
			}
		}
	}
	if input.Err() != nil {
		log.Logger.WithField("remoteAddr", c.rawConn.RemoteAddr().String()).Error(input.Err())
		if c.getState() == connectSt {
			// disconnection by the other side
			c.setState(remoteDisConnSt)
			c.disConnCh <- c.peerAddr
		}
		return
	}
}

func (c *Conn) specialMsg(msg []byte) bool {
	switch msgType := msg[0]; msgType {
	case handshakeCode:
		b, err := json.Marshal(c.ourNodeInfo)
		if err != nil {
			log.Logger.Errorf("special message json marshal error:%v", err)
		}
		if err = c.write(bytes.Join([][]byte{
			{handshakeRespCode},
			b,
		}, []byte{})); err != nil {
			log.Logger.Errorf("special message write error:%v", err)
		}
		return true
	}
	return false
}

func (c *Conn) write(msg []byte) error {
	return write(c.rawConn, msg)
}

func (c *Conn) getState() connState {
	defer c.state.mtx.RUnlock()
	c.state.mtx.RLock()
	return c.state.val
}

func (c *Conn) setState(a connState) {
	defer c.state.mtx.Unlock()
	c.state.mtx.Lock()
	c.state.val = a
}

func write(rawConn net.Conn, msg []byte) error {
	w := bufio.NewWriter(rawConn)
	_, err := w.Write(bytes.Join([][]byte{
		msg,
		[]byte("\n"),
	}, []byte{}))
	if err != nil {
		return err
	}
	return w.Flush()
}

// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package p2p

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/reed/errors"
	"github.com/reed/log"
	"github.com/tendermint/tmlibs/common"
	"net"
	"sync"
	"time"
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

func NewConnection(peerListenAddr string, disConnCh chan<- string, rawConn net.Conn, ourNodeInfo *NodeInfo, handleFunc HandleFunc) *Conn {
	rawConn.SetReadDeadline(time.Time{}) // reset read deadline:not time out.
	conn := &Conn{
		ourNodeInfo: ourNodeInfo,
		peerAddr:    peerListenAddr,
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
	} else {
		log.Logger.Info("closed conn")
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
				log.Logger.Errorf("connection.read failed to write:%v", err)
			}
		}
	}
	if input.Err() != nil {
		log.Logger.WithField("remoteAddr", c.rawConn.RemoteAddr().String()).Errorf("readGoroutine error:%v", input.Err())
	} else {
		log.Logger.WithField("remoteAddr", c.rawConn.RemoteAddr().String()).Errorf("readGoroutine has closed")
	}
	if c.getState() == connectSt {
		// disconnection by the other side
		c.setState(remoteDisConnSt)
		c.disConnCh <- c.peerAddr
	}
	return
}

func (c *Conn) specialMsg(msg []byte) bool {
	switch msgType := msg[0]; msgType {
	case handshakeCode:
		if err := writeOurNodeInfo(c.rawConn, c.ourNodeInfo); err != nil {
			log.Logger.Error(err)
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

func writeOurNodeInfo(rawConn net.Conn, ourNodeInfo *NodeInfo) error {
	b, err := json.Marshal(ourNodeInfo)
	if err != nil {
		return errors.Wrap(err, "special message json marshal error")
	}
	if err = write(rawConn, bytes.Join([][]byte{
		{handshakeRespCode},
		b,
	}, []byte{})); err != nil {
		return errors.Wrap(err, "special message write error")
	}
	return nil
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

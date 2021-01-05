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

type Peer struct {
	common.BaseService
	nodeInfo *NodeInfo
	conn     *Conn
}

func NewPeer(ourNodeInfo *NodeInfo, nodeInfo *NodeInfo, disConnCh chan<- string, rawConn net.Conn, handleFunc HandleFunc) *Peer {
	peer := &Peer{
		nodeInfo: nodeInfo,
		conn:     NewConnection(nodeInfo.RemoteAddr, disConnCh, rawConn, ourNodeInfo, handleFunc),
	}
	peer.BaseService = *common.NewBaseService(nil, "peer", peer)
	return peer
}

func (p *Peer) OnStart() error {
	return p.conn.Start()
}

func (p *Peer) OnStop() {
	if err := p.conn.Stop(); err != nil {
		log.Logger.Error(err)
	}
}

func toAddress(ip net.IP, port uint16) string {
	return net.JoinHostPort(
		ip.String(),
		strconv.FormatUint(uint64(port), 10),
	)
}

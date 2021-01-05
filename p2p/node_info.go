// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package p2p

import (
	"encoding/json"
	"github.com/reed/errors"
	"github.com/reed/p2p/discover"
	"net"
	"strconv"
)

var (
	newNodeInfoErr = errors.New("new nodeInfo error")
)

type NodeInfo struct {
	ID         discover.NodeID
	RemoteAddr string
	// Version    uint64
}

func NewOurNodeInfo(id discover.NodeID, ip net.IP, port uint16) *NodeInfo {
	return &NodeInfo{
		ID:         id,
		RemoteAddr: net.JoinHostPort(ip.String(), strconv.FormatUint(uint64(port), 10)),
	}
}

func NewNodeInfoFromBytes(b []byte, conn net.Conn) (*NodeInfo, error) {
	ni := &NodeInfo{}
	if err := json.Unmarshal(b, ni); err != nil {
		return nil, errors.Wrap(newNodeInfoErr, err)
	}
	ni.RemoteAddr = conn.RemoteAddr().String()
	return ni, nil
}

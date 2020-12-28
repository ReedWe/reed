// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package discover

import (
	"encoding/hex"
	"github.com/reed/blockchain/config"
	"github.com/reed/log"
	"net"
	"strconv"
	"strings"
)

func getSeeds() []*Node {
	var nodes []*Node
	seeds := config.Default.Seeds
	if seeds == "" {
		log.Logger.Info("no seeds")
		return nodes
	}

	arr := strings.Split(seeds, ",")
	for _, a := range arr {
		if n := resolveNode(a); n != nil {
			nodes = append(nodes, n)
		}
	}
	return nodes
}

func resolveNode(enode string) *Node {
	if enode == "" {
		return nil
	}

	e := strings.Split(enode, "@")

	id, err := hex.DecodeString(e[0])
	if err != nil {
		log.Logger.Errorf("failed to resolve enode:%v", err)
		return nil
	}

	p := strings.Split(e[1], ":")
	addr, _ := net.ResolveIPAddr("ip", p[0])

	parseUint, err := strconv.ParseUint(p[1], 10, 16)
	if err != nil {
		log.Logger.Errorf("failed to resolve enode:%v", err)
		return nil
	}

	return NewNode(BytesToHash(id), addr.IP, uint16(parseUint))
}

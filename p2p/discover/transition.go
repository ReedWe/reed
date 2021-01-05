// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package discover

import (
	"github.com/reed/errors"
	"github.com/reed/log"
	"github.com/sirupsen/logrus"
)

var (
	invalidNodeEvent = errors.New("invalid node event")
)

var (
	unknown          *nodeState
	verifyInit       *nodeState
	verifyWait       *nodeState
	remoteVerifyWait *nodeState
	known            *nodeState
	contested        *nodeState
	unresponsive     *nodeState
)

type nodeState struct {
	name   string
	enter  func(u *UDP, n *Node)
	handle func(u *UDP, n *Node, e nodeEvent, pkt *ingressPacket) (*nodeState, error)
}

func init() {
	unknown = &nodeState{
		name: "unknown",
		enter: func(u *UDP, n *Node) {
			u.table.delete(n.ID)

			// reset timeout count
			n.queryTimeouts = 0

			// clear defer query
			for _, q := range n.deferQueries {
				if q.reply != nil {
					q.reply <- nil
				}
			}
			n.deferQueries = nil

			// clear pending query
			if n.pendingQuery != nil {
				if n.pendingQuery != nil {
					n.pendingQuery.reply <- nil
				}
				n.pendingQuery = nil
			}
		},
		handle: func(u *UDP, n *Node, e nodeEvent, pkt *ingressPacket) (*nodeState, error) {
			log.Logger.WithFields(logrus.Fields{"remoteIP": n.IP, "Port": n.UDPPort, "event": e, "nodeID": n.ID.ToString()}).Debug("transition[unknown].handle")
			switch e {
			case pingPacket:
				// response immediately
				u.sendPong(pkt)
				// ping for ensure that the node is reliable
				u.sendPing(n)
				return verifyWait, nil
			default:
				return unknown, invalidNodeEvent
			}
		},
	}

	verifyInit = &nodeState{
		name: "verifyInit",
		enter: func(u *UDP, n *Node) {
			log.Logger.Infof("verifyInit enter")
			u.sendPing(n)
		},
		handle: func(u *UDP, n *Node, e nodeEvent, pkt *ingressPacket) (*nodeState, error) {
			log.Logger.WithFields(logrus.Fields{"remoteIP": n.IP, "Port": n.UDPPort, "event": e, "nodeID": n.ID.ToString()}).Debug("transition[verifyInit].handle")
			switch e {
			case pingPacket:
				u.sendPong(pkt)
				return verifyInit, nil
			case pongPacket:
				u.processPong(n)
				return remoteVerifyWait, nil
			case pongTimeout:
				return unknown, nil
			default:
				return verifyInit, invalidNodeEvent
			}
		},
	}

	verifyWait = &nodeState{
		name: "verifyWait",
		handle: func(u *UDP, n *Node, e nodeEvent, pkt *ingressPacket) (*nodeState, error) {
			log.Logger.WithFields(logrus.Fields{"remoteIP": n.IP, "Port": n.UDPPort, "event": e, "nodeID": n.ID.ToString()}).Debug("transition[verifyWait].handle")
			switch e {
			case pingPacket:
				u.sendPong(pkt)
				return verifyWait, nil
			case pongPacket:
				u.processPong(n)
				return known, nil
			case pongTimeout:
				return unknown, nil
			default:
				return verifyWait, invalidNodeEvent
			}
		},
	}

	remoteVerifyWait = &nodeState{
		name: "remoteVerifyWait",
		enter: func(u *UDP, n *Node) {
			u.makeTimeoutEvent(n, pingTimeout, responseTimeout)
		},
		handle: func(u *UDP, n *Node, e nodeEvent, pkt *ingressPacket) (*nodeState, error) {
			log.Logger.WithFields(logrus.Fields{"remoteIP": n.IP, "Port": n.UDPPort, "event": e, "nodeID": n.ID.ToString()}).Debug("transition[remoteVerifyWait].handle")
			switch e {
			case pingPacket:
				u.sendPong(pkt)
				return remoteVerifyWait, nil
			case pingTimeout:
				return known, nil
			default:
				return remoteVerifyWait, invalidNodeEvent
			}
		},
	}

	known = &nodeState{
		name: "know",
		enter: func(u *UDP, n *Node) {
			log.Logger.WithFields(logrus.Fields{"remoteIP": n.IP, "Port": n.UDPPort, "nodeID": n.ID.ToString()}).Debug("=know node=")
			n.queryTimeouts = 0
			n.executeDefer(u)
			u.table.add(n)
		},
		handle: func(u *UDP, n *Node, e nodeEvent, pkt *ingressPacket) (*nodeState, error) {
			log.Logger.WithFields(logrus.Fields{"remoteIP": n.IP, "Port": n.UDPPort, "event": e, "nodeID": n.ID.ToString()}).Debug("transition[know].handle")
			switch e {
			case pingPacket:
				u.sendPong(pkt)
				return known, nil
			case pongPacket:
				u.processPong(n)
				return known, nil
			default:
				return u.processQueryEvent(n, e, pkt)
			}
		},
	}

	contested = &nodeState{
		name: "contested",
		enter: func(u *UDP, n *Node) {
			u.sendPing(n)
		},
		handle: func(u *UDP, n *Node, e nodeEvent, pkt *ingressPacket) (*nodeState, error) {
			log.Logger.WithFields(logrus.Fields{"remoteIP": n.IP, "Port": n.UDPPort, "event": e, "nodeID": n.ID.ToString()}).Debug("transition[contested].handle")
			switch e {
			case pingPacket:
				u.sendPong(pkt)
				return contested, nil
			case pongPacket:
				u.processPong(n)
				return known, nil
			case pongTimeout:
				return unresponsive, nil
			default:
				return u.processQueryEvent(n, e, pkt)
			}
		},
	}

	unresponsive = &nodeState{
		name: "unresponsive",
		handle: func(u *UDP, n *Node, e nodeEvent, pkt *ingressPacket) (*nodeState, error) {
			log.Logger.WithFields(logrus.Fields{"remoteIP": n.IP, "Port": n.UDPPort, "event": e, "nodeID": n.ID.ToString()}).Debug("transition[unresponsive].handle")
			switch e {
			case pingPacket:
				u.sendPong(pkt)
				return known, nil
			case pongPacket:
				u.processPong(n)
				return known, nil
			default:
				return u.processQueryEvent(n, e, pkt)
			}
		},
	}
}

func transform(u *UDP, n *Node, nextState *nodeState) {
	log.Logger.WithFields(logrus.Fields{"remoteIP": n.IP, "Port": n.UDPPort, "state": n.state, "nextState": nextState}).Info("node state transform")
	if n.state != nextState {
		n.state = nextState
		if nextState.enter != nil {
			nextState.enter(u, n)
		}
	}
}

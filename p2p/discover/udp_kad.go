// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package discover

import (
	"github.com/reed/log"
	"github.com/sirupsen/logrus"
)

func (g *nodeNetGuts) pushToDefer(f *findNodeQuery) {
	g.deferQueries = append(g.deferQueries, f)
}

func (g *nodeNetGuts) executeDefer(u *UDP) {
	if len(g.deferQueries) == 0 {
		log.Logger.Info("no query request in deferQueries")
		return
	}
	q := g.deferQueries[0]
	if q.maybeExecute(u) {
		// move on
		g.deferQueries = append(g.deferQueries[:0], g.deferQueries[1:]...)
	}
}

func (f *findNodeQuery) maybeExecute(u *UDP) (success bool) {
	log.Logger.WithFields(logrus.Fields{"remoteID": f.remote.ID.ToString(), "IP": f.remote.IP, "Port": f.remote.UDPPort, "target": f.target.ToString()}).Info("try execute findNodeQuery")
	if f.remote == u.OurNode {
		// satisfy queries against the local node directly.
		log.Logger.Debug("remote is self,return the local node directly")
		nd := u.table.closest(f.target)
		log.Logger.Infof("--closest node count:%d", len(nd.entries))
		for _, entry := range nd.entries {
			log.Logger.WithFields(logrus.Fields{"IP": entry.IP, "Port": entry.UDPPort, "ID": entry.ID.ToString()}).Debug("--closest")
		}
		if f.reply != nil {
			f.reply <- &findNodeRespReply{remoteID: f.remote.ID, nodes: nd.entries}
		}
		return
	}

	if f.remote.state == unknown {
		log.Logger.Info("unknown node")
		transform(u, f.remote, verifyInit)
		return
	}
	if !f.remote.state.canQuery() {
		log.Logger.WithField("state", f.remote.state).Info("can not query:invalid state")
		return
	}
	if f.remote.pendingQuery != nil {
		log.Logger.Info("can not query:pendingQuery Queue not empty")
		return
	}
	u.sendFindNode(f.remote, f.target)
	f.remote.pendingQuery = f
	return true
}

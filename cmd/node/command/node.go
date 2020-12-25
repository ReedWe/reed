// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package command

import (
	"github.com/prometheus/tsdb/fileutil"
	"github.com/reed/api"
	bc "github.com/reed/blockchain"
	"github.com/reed/blockchain/config"
	"github.com/reed/blockchain/store"
	"github.com/reed/database/leveldb"
	"github.com/reed/errors"
	"github.com/reed/log"
	"github.com/reed/miner"
	"github.com/reed/p2p"
	"github.com/reed/wallet"
	"github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
	"path/filepath"
)

type Node struct {
	common.BaseService
	api          *api.API
	chain        *bc.Chain
	miner        *miner.Miner
	instanceLock fileutil.Releaser
	discover     *p2p.Server
}

func NewNode() *Node {
	releaser, err := lockDataDir()
	if err != nil {
		common.Exit(err.Error())
	}

	s := newStore()
	chain, err := bc.NewChain(&s)
	if err != nil {
		common.Exit(common.Fmt("Failed to create chain:%v", err))
	}

	w, _ := wallet.My("123")
	node := &Node{
		api:          api.NewApi(chain),
		chain:        chain,
		miner:        miner.NewMiner(chain, w, chain.GetWriteReceptionChan(), chain.GetReadBreakWorkChan()),
		instanceLock: releaser,
	}

	node.BaseService = *common.NewBaseService(nil, "Node", node)

	p2p, err := p2p.NewP2PServer()
	if err != nil {
		common.Exit(common.Fmt(err.Error()))
	}
	node.discover = p2p

	return node
}

func (n *Node) OnStart() error {
	n.api.StartApiServer()
	if err := n.chain.Open(); err != nil {
		return err
	}
	if config.Default.Mining {
		if err := n.miner.Start(); err != nil {
			return err
		}
	}
	n.discover.Start()
	log.Logger.Info("Node started successfully.")
	return nil
}

func (n *Node) OnStop() {
	n.chain.Close()
	n.miner.Stop()
	if err := n.instanceLock.Release(); err != nil {
		log.Logger.Errorf("Can't release dataDir locke:%v", err)
	}
	n.instanceLock = nil
	log.Logger.Info("Node has shut down.")
}

func (n *Node) RunForever() {
	common.TrapSignal(func() {
		n.Stop()
	})
}

func lockDataDir() (fileutil.Releaser, error) {
	lock, _, err := fileutil.Flock(filepath.Join(config.Default.HomeDir, config.Default.LockName))
	if err != nil {
		return nil, errors.Wrapf(err, "Can not start node")
	}
	return lock, nil
}

func newStore() store.Store {
	return leveldb.NewStore(dbm.NewDB("core", dbm.LevelDBBackend, config.DatabaseDir()))
}

// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package command

import (
	"github.com/reed/api"
	bc "github.com/reed/blockchain"
	"github.com/reed/blockchain/store"
	"github.com/reed/database/leveldb"
	"github.com/reed/log"
	"github.com/reed/miner"
	"github.com/reed/wallet"
	"github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
	"os"
)

type Node struct {
	common.BaseService
	api   *api.API
	chain *bc.Chain
	miner *miner.Miner
}

func NewNode() *Node {
	s := newStore()
	chain, err := bc.NewChain(&s)
	if err != nil {
		common.Exit(common.Fmt("Failed to create chain:%v", err))
	}

	w, _ := wallet.My("123")
	node := &Node{
		api:   api.NewApi(chain),
		chain: chain,
		miner: miner.NewMiner(chain, w, chain.GetWriteReceptionChan(), chain.GetReadBreakWorkChan()),
	}

	node.BaseService = *common.NewBaseService(nil, "Node", node)

	return node
}

func (n *Node) OnStart() error {
	n.api.StartApiServer()
	if err := n.chain.Open(); err != nil {
		return err
	}
	if err := n.miner.Start(); err != nil {
		return err
	}
	log.Logger.Info("Node started successfully.")
	return nil
}

func (n *Node) OnStop() {
	n.chain.Close()
	n.miner.Stop()
	log.Logger.Info("Node has shut down.")
}

func (n *Node) RunFover() {
	common.TrapSignal(func() {

	})
}

func newStore() store.Store {
	return leveldb.NewStore(dbm.NewDB("core", dbm.LevelDBBackend, os.Getenv("GOPATH")+"/src/github.com/reed/database/file/"))
}

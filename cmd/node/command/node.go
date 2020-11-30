package command

import (
	"github.com/reed/api"
	bc "github.com/reed/blockchain"
	"github.com/reed/database/leveldb"
	"github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
	"os"
)

type Node struct {
	common.BaseService
	api   *api.API
	chain *bc.Chain
}

func NewNode() *Node {

	leveldbStore := leveldb.NewStore(dbm.NewDB("core", dbm.LevelDBBackend, os.Getenv("GOPATH")+"/src/github.com/reed/database/file/"))

	chain, err := bc.NewChain(leveldbStore)
	if err != nil {
		common.Exit(common.Fmt("Failed to create chain:%v", err))
	}

	node := &Node{
		api:   api.NewApi(chain),
		chain: chain,
	}

	node.BaseService = *common.NewBaseService(nil, "Node", node)

	return node
}

func (n *Node) OnStart() error {
	n.api.StartApiServer()
	return nil
}

func (n *Node) OnStop() {
	//do nothing
}

func (n *Node) RunFover() {
	common.TrapSignal(func() {

	})
}

package command

import (
	"github.com/tendermint/tmlibs/common"
	"github.com/reed/api"
)

type Node struct {
	common.BaseService
	api *api.API
}

func NewNode() *Node {
	node := &Node{
		api: api.NewApi(),
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

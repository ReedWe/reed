package node

import (
	"fmt"
	"github.com/tendermint/tmlibs/common"
	"github.com/tybc/api"
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
	fmt.Println("Node Onstart...")

	n.api.StartApiServer()
	return nil
}

func (n *Node) OnStop() {
	fmt.Println("Node Onstop")
}

func (n *Node) RunFover() {
	common.TrapSignal(func() {

	})
}

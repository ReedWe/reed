package blockchain

import "github.com/tybc/core"

type Chain struct {
	Store  Store
	Txpool *core.Txpool
}

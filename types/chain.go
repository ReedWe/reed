package types

import (
	"github.com/tybc/blockchain/store"
	"github.com/tybc/blockchain/txpool"
)

type Chain struct {
	Store  store.Store
	Txpool *txpool.Txpool
}

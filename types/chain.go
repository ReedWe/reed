package types

import (
	"github.com/reed/blockchain/store"
	"github.com/reed/blockchain/txpool"
)

type Chain struct {
	Store  store.Store
	Txpool *txpool.Txpool
}

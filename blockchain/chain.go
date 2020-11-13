package blockchain

import "github.com/tybc/blockchain/txpool"

type Chain struct {
	Store  Store
	Txpool *txpool.Txpool
}

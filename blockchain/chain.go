package blockchain

import (
	bm "github.com/reed/blockchain/blockmanager"
	"github.com/reed/blockchain/store"
	"github.com/reed/blockchain/txpool"
)

type Chain struct {
	Store        store.Store
	Txpool       *txpool.Txpool
	BlockManager *bm.BlockManager
}

func NewChain(s store.Store) (*Chain, error) {
	tp := txpool.NewTxpool(&s)
	highestBlock, err := s.GetHighestBlock()
	if err != nil {
		return nil, err
	}
	blockMgr, err := bm.NewBlockManager(&s, highestBlock)
	if err != nil {
		return nil, err
	}
	return &Chain{
		Store:        s,
		Txpool:       tp,
		BlockManager: blockMgr,
	}, nil
}

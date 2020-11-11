package blockchain

import (
	"github.com/tybc/core/types"
	"sync"
)

type Txpool struct {
	Txs       map[types.Hash]*types.Tx
	OutputIds map[types.Hash]*types.TxOutput
	Store     Store
	mtx       sync.RWMutex
}

func NewTxpool(store Store) *Txpool {
	return &Txpool{
		Txs:   make(map[types.Hash]*types.Tx),
		Store: store,
	}
}


func (tp *Txpool) ExistOutput(hash types.Hash) bool {
	tp.mtx.RLock()
	defer tp.mtx.RUnlock()

	return tp.OutputIds[hash] != nil
}

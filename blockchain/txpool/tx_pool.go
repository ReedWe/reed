package txpool

import (
	"github.com/tybc/blockchain"
	"github.com/tybc/types"
	"sync"
)

type Txpool struct {
	Txs       map[types.Hash]*types.Tx
	OutputIds map[types.Hash]*types.TxOutput
	Store     blockchain.Store
	mtx       sync.RWMutex
}

func NewTxpool(store blockchain.Store) *Txpool {
	return &Txpool{
		Txs:   make(map[types.Hash]*types.Tx),
		Store: store,
	}
}

func (tp *Txpool) GetTx(txId *types.Hash) (*types.Tx, error) {
	tp.mtx.RLock()
	defer tp.mtx.RUnlock()

	return tp.Store.GetTx((*txId).Bytes())
}

func (tp *Txpool) ExistOutput(hash types.Hash) bool {
	tp.mtx.RLock()
	defer tp.mtx.RUnlock()

	return tp.OutputIds[hash] != nil
}

package txpool

import (
	"github.com/reed/blockchain/store"
	"github.com/reed/types"
	"sync"
)

type Txpool struct {
	Txs       map[types.Hash]*types.Tx
	OutputIds map[types.Hash]*types.TxOutput
	Store     store.Store
	mtx       sync.RWMutex
}

func NewTxpool(store *store.Store) *Txpool {
	return &Txpool{
		Txs:   make(map[types.Hash]*types.Tx),
		Store: *store,
	}
}

func (tp *Txpool) AddTx(tx *types.Tx) error {
	tp.mtx.Lock()
	defer tp.mtx.Unlock()

	return tp.Store.AddTx(tx)
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

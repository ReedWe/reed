// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package txpool

import (
	"github.com/reed/blockchain/store"
	"github.com/reed/errors"
	"github.com/reed/types"
	"sync"
)

type Txpool struct {
	Txs   map[types.Hash]*types.Tx
	Store store.Store
	mtx   sync.RWMutex
}

var (
	addTxErr = errors.New("transaction pool add error")
)

func NewTxpool(store *store.Store) *Txpool {
	return &Txpool{
		Txs:   map[types.Hash]*types.Tx{},
		Store: *store,
	}
}

func (tp *Txpool) AddTx(tx *types.Tx) error {
	tp.mtx.Lock()
	defer tp.mtx.Unlock()

	if _, ok := tp.Txs[tx.GetID()]; ok {
		return errors.Wrapf(addTxErr, "tx exists ID=%x", tx.GetID())
	}
	tp.Txs[tx.GetID()] = tx
	return nil
}

func (tp *Txpool) GetTx(txId types.Hash) *types.Tx {
	tp.mtx.RLock()
	defer tp.mtx.RUnlock()

	return tp.Txs[txId]
}

func (tp *Txpool) GetTxs() []*types.Tx {
	tp.mtx.RLock()
	defer tp.mtx.RUnlock()

	txs := make([]*types.Tx, len(tp.Txs), len(tp.Txs)+1)
	for _, t := range tp.Txs {
		txs = append(txs, t)
	}
	return txs
}

func (tp *Txpool) RemoveTransactions(txs []*types.Tx) {
	tp.mtx.Lock()
	defer tp.mtx.Unlock()

	for _, tx := range txs {
		delete(tp.Txs, tx.GetID())
	}
}

// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package store

import (
	"github.com/reed/types"
)

type Store interface {
	// SaveTx(tx *types.Tx) error
	//
	// GetTx(id []byte) (*types.Tx, error)

	GetUtxo(id []byte) (*types.UTXO, error)

	SaveUtxos(usedUtxoIDs []*types.Hash, utxos []*types.UTXO) error

	GetHighestBlock() (*types.Block, error)

	GetBlock(hash []byte) (*types.Block, error)

	SaveBlock(block *types.Block) error
}

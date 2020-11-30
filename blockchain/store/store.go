package store

import "github.com/reed/types"

type Store interface {
	AddTx(tx *types.Tx) error

	GetTx(id []byte) (*types.Tx, error)

	GetUtxo(id []byte) (*types.UTXO, error)

	SaveUtxos(expiredUtxoIds []types.Hash, utxos *[]types.UTXO) error

	GetHighestBlock() (*types.Hash, error)

	GetBlock(hash []byte) (*types.Block, error)

	AddBlock(block *types.Block) error
}

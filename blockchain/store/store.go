package store

import "github.com/reed/types"

type Store interface {
	SaveTx(tx *types.Tx) error

	GetTx(id []byte) (*types.Tx, error)

	GetUtxo(id []byte) (*types.UTXO, error)

	SaveUtxos(expiredUtxoIds []types.Hash, utxos *[]types.UTXO) error

	GetHighestBlock() (*types.Block, error)

	GetBlock(hash []byte) (*types.Block, error)

	SaveBlockAndUpdateHighest(block *types.Block) error
}

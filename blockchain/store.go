package blockchain

import "github.com/tybc/types"

type Store interface {
	AddTx(tx *types.Tx) error

	GetTx(id []byte) (*types.Tx, error)

	GetUtxo(id []byte) (*types.UTXO, error)

	SaveUtxo(utxo *types.TxOutput)
}

package blockchain

import "github.com/tybc/types"

type Store interface {
	GetUtxo(id []byte) (*types.UTXO, error)
	SaveUtxo(utxo *types.TxOutput)
}

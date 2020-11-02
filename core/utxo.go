package core

import "github.com/tybc/core/types"

type Utxo struct {
	types.TxOutput
}

func GetUtxoById(id []byte) *Utxo {
	return &Utxo{}
}

package types

import (
	"github.com/tybc/blockchain"
	"github.com/tybc/core"
)

type TxInput struct {
	SpendOutputId Hash
	*UTXO

	ScriptSig []byte
}

func (txInput *TxInput) SetUtxo(store *blockchain.Store) error {
	if utxo, err := core.GetUtxoByOutputId(store, txInput.SpendOutputId); err != nil {
		return err
	} else {
		txInput.UTXO = utxo
	}
	return nil
}

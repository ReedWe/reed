package blockchain

import (
	"github.com/tybc/errors"
	"github.com/tybc/types"
)

func GetUtxoByOutputId(store *Store, outputId types.Hash) (*types.UTXO, error) {
	id := outputId.Bytes()
	utxo, err := (*store).GetUtxo(id)
	if utxo == nil || err != nil {
		return nil, errors.Wrapf(err, "utxo(outputId:%x) does not exists", id)
	}
	return utxo, nil
}

func OutputsToUtxos(txId *types.Hash, outputs *[]types.TxOutput) *[]types.UTXO {
	var utxos []types.UTXO
	for i, output := range *outputs {
		utxos = append(utxos, types.UTXO{
			OutputId:  output.ID,
			SoureId:   *txId,
			SourcePos: uint32(i),
			Amount:    output.Amount,
			Address:   output.Address,
			ScriptPk:  output.ScriptPk,
		})
	}
	return &utxos
}

func UtxoChange(store *Store, inputs *[]types.TxInput, utxos *[]types.UTXO) error {
	var expiredUtxoIds []types.Hash
	for _, input := range *inputs {
		expiredUtxoIds = append(expiredUtxoIds, input.ID)
	}

	return (*store).SaveUtxos(expiredUtxoIds, utxos)
}

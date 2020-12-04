// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package blockchain

import (
	"github.com/reed/blockchain/store"
	"github.com/reed/errors"
	"github.com/reed/types"
)

func ProcessUtxoForSaveBlock(store *store.Store, b *types.Block) error {
	var utxos []*types.UTXO
	var usedUtxoIDs []*types.Hash

	for _, tx := range b.Transactions {
		txID := tx.GetID()
		utxos = append(utxos, OutputsToUtxos(&txID, tx.TxOutput)...)
		usedUtxoIDs = append(usedUtxoIDs, InputsToUtxoIDs(tx.TxInput)...)
	}
	return UtxoChange(store, usedUtxoIDs, utxos)
}

func GetUtxoByOutputId(store *store.Store, outputId types.Hash) (*types.UTXO, error) {
	id := outputId.Bytes()
	utxo, err := (*store).GetUtxo(id)
	if utxo == nil || err != nil {
		return nil, errors.Wrapf(err, "utxo(outputId:%x) does not exists", id)
	}
	return utxo, nil
}

func OutputsToUtxos(txId *types.Hash, outputs []*types.TxOutput) []*types.UTXO {
	var utxos []*types.UTXO
	for i, output := range outputs {
		utxos = append(utxos,
			types.NewUtxo(output.ID, *txId, output.IsCoinBase, uint64(i), output.Amount, output.Address, output.ScriptPk),
		)
	}
	return utxos
}

func InputsToUtxoIDs(inputs []*types.TxInput) []*types.Hash {
	var IDs []*types.Hash
	for _, input := range inputs {
		IDs = append(IDs, &input.SpendOutputId)
	}
	return IDs
}

func UtxoChange(store *store.Store, usedUtxoIDs []*types.Hash, utxos []*types.UTXO) error {
	return (*store).SaveUtxos(usedUtxoIDs, utxos)
}

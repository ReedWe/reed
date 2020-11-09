package validation

import (
	"github.com/tybc/blockchain"
	"github.com/tybc/core/types"
	"github.com/tybc/errors"
	"github.com/tybc/log"
)

func CheckUtxoExists(store *blockchain.Store, tx *types.Tx) error {
	for _, input := range tx.TxInput {
		spendOutId := input.SpendOutputId[:]
		log.Logger.Infof("utxo spendOutputId:%x", spendOutId)
		utxo, err := (*store).GetUtxo(spendOutId)
		if utxo == nil || err != nil {
			return errors.Wrapf(err, "utxo(spendOutputId:%x) does not exists", spendOutId)
		}
	}
	return nil
}

//func Check(tx *types.Transaction) error {
//	err := CheckUtxoExists(tx.TxOutput)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

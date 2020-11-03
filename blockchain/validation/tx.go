package validation

import (
	"github.com/tybc/blockchain"
	"github.com/tybc/errors"
	"github.com/tybc/log"
)

func CheckUtxoExists(store *blockchain.Store, ids *[][]byte) error {
	for _, id := range *ids {
		log.Logger.Infof("utxo id:%x", id)
		utxo, err := (*store).GetUtxo(id)
		if utxo == nil || err != nil {
			return errors.Wrapf(err, "utxo(id:%x) does not exists", id)
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

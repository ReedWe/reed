// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package txpusher

import (
	"github.com/reed/blockchain"
	"github.com/reed/blockchain/validation"
	"github.com/reed/errors"
	"github.com/reed/log"
	"github.com/reed/types"
)

var (
	recvTxErr = errors.New("receive transaction error")
)

// receive local transaction and remote transaction
func MaybePush(chain *blockchain.Chain, tx *types.Tx) error {
	log.Logger.Infof("receive a new transaction ID=%x", tx.ID)

	getUtxo := func(spendOutputId types.Hash) (*types.UTXO, error) {
		return blockchain.GetUtxoByOutputId(chain.Store, spendOutputId)
	}
	if err := tx.Completion(getUtxo); err != nil {
		return err
	}

	//generate txID and validate
	txId, err := tx.GenerateID()
	if !txId.HashEqual(tx.ID) {
		return errors.Wrapf(recvTxErr, "txId errors. local=%x remote=%x", txId, tx.ID)
	}

	existTx := chain.Txpool.GetTx(tx.ID)

	if existTx != nil {
		log.Logger.Infof("transaction exists already (id=%x)", tx.ID)
		return nil
	}

	if err := validation.ValidateTx(tx); err != nil {
		return err
	}

	//push into tx pool
	if err = chain.Txpool.AddTx(tx); err != nil {
		return err
	}

	//TODO broadcast this transaction

	// TODO NOT HERE
	//output utxo and save changed
	utxos := blockchain.OutputsToUtxos(&tx.ID, tx.TxOutput)
	if err = blockchain.UtxoChange(chain.Store, tx.TxInput, utxos); err != nil {
		return err
	}

	return nil
}

// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package types

import (
	"github.com/reed/common/math"
	"github.com/reed/crypto"
	"github.com/reed/errors"
	"strconv"
)

var (
	txOutputCheckErr = errors.New("transaction output check error")
)

type Tx struct {
	ID Hash `json:"id"`

	TxInput  []*TxInput  `json:"input"`
	TxOutput []*TxOutput `json:"output"`
}

func (tx *Tx) GenerateID() (*Hash, error) {
	if len(tx.TxInput) == 0 {
		return nil, errors.New("tx input empty")
	}

	var ids [][]byte
	for _, input := range tx.TxInput {
		ids = append(ids, input.ID.Bytes())
	}
	for _, output := range tx.TxOutput {
		ids = append(ids, output.ID.Bytes())
	}

	h := BytesToHash(crypto.Sha256(ids...))
	return &h, nil
}

func (tx *Tx) Completion(getUtxo func(spendOutputId Hash) (*UTXO, error)) error {
	//set spend data
	//check input rel utxo
	for _, input := range tx.TxInput {
		if utxo, err := getUtxo(input.SpendOutputId); err != nil {
			return err
		} else {
			input.SetSpend(utxo)
		}
	}

	//input ID
	for _, input := range tx.TxInput {
		//ID
		input.ID = input.GenerateID()
	}

	//output ID
	//locking script
	for _, output := range tx.TxOutput {
		if output.Amount == 0 {
			return errors.Wrapf(txOutputCheckErr, "invalid output amount:%d", output.Amount)
		}
		if len(output.Address) != 32 {
			return errors.Wrapf(txOutputCheckErr, "invalid output address. len(%d).expect 36", len(output.Address))
		}
	}
	return nil
}

func (tx *Tx) IsAssetAmtEqual() (sumInput uint64, sumOutput uint64, err error) {
	for _, input := range tx.TxInput {
		if sumInput, err = math.AddUint64(sumInput, input.Amount); err != nil {
			return 0, 0, err
		}
	}
	for _, output := range tx.TxOutput {
		if sumOutput, err = math.AddUint64(sumOutput, output.Amount); err != nil {
			return 0, 0, err
		}
	}
	return
}

func NewCoinbaseTx(curHeight uint64, coinbaseAddr []byte, amt uint64) (*Tx, error) {
	i := &TxInput{ScriptSig: []byte("height:" + strconv.FormatUint(curHeight, 10))}
	i.ID = i.GenerateID()

	o := NewTxOutput(true, coinbaseAddr, amt)

	var inps []*TxInput
	var iops []*TxOutput

	inps = append(inps, i)
	iops = append(iops, o)
	tx := &Tx{TxInput: inps, TxOutput: iops}
	id, err := tx.GenerateID()
	if err != nil {
		return nil, err
	}
	tx.ID = *id
	return tx, nil
}

//func (tx *Transaction) sign() {
//
//	for _, input := range tx.TxInput {
//		input.Signature = crypto.Sha256(input,tx.SetID)
//	}
//
//}

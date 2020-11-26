package types

import (
	"github.com/reed/common/math"
	"github.com/reed/crypto"
	"github.com/reed/errors"
)

var (
	txOutputCheckErr = errors.New("transaction outpu check error")
)

type Tx struct {
	ID Hash `json:"id"`

	TxInput  []TxInput  `json:"input"`
	TxOutput []TxOutput `json:"output"`
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
		output.ScriptPk = output.GenerateLockingScript()
		output.ID = output.GenerateID()
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

//func (tx *Transaction) sign() {
//
//	for _, input := range tx.TxInput {
//		input.Signature = crypto.Sha256(input,tx.SetID)
//	}
//
//}

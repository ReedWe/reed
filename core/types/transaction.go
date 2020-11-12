package types

import (
	"github.com/tybc/crypto"
	"github.com/tybc/errors"
)

//id = Hash(input1.id,input2.id,...)
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

//func (tx *Transaction) sign() {
//
//	for _, input := range tx.TxInput {
//		input.Signature = crypto.Sha256(input,tx.SetID)
//	}
//
//}

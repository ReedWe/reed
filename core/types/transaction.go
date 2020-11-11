package types

import "github.com/tybc/crypto"

//id = Hash(input1.id,input2.id,...)
type Tx struct {
	ID Hash `json:"id"`

	TxInput  []TxInput  `json:"input"`
	TxOutput []TxOutput `json:"output"`
}

func (tx *Tx) SetID() error {
	ids := make([][]byte, len(tx.TxInput))
	for _, input := range tx.TxInput {
		ids = append(ids, input.ID[:])
	}
	tx.ID = BytesToHash(crypto.Sha256(ids...))
	return nil
}

//func (tx *Transaction) sign() {
//
//	for _, input := range tx.TxInput {
//		input.Signature = crypto.Sha256(input,tx.SetID)
//	}
//
//}

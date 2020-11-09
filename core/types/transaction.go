package types

type Tx struct {
	ID [32]byte

	TxInput  []TxInput
	TxOutput []TxOutput

	*UTXO
}

//func (tx *Transaction) sign() {
//
//	for _, input := range tx.TxInput {
//		input.Signature = crypto.Sha256(input,tx.ID)
//	}
//
//}

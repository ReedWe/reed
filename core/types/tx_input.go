package types

import (
	"bytes"
	"github.com/tybc/crypto"
	"github.com/tybc/errors"
)

var (
	inputErr = errors.New("tx input")
)

type TxInput struct {
	Spend
	ScriptSig []byte //Sig(Hash(TxInput.Id + tx.Id))
}

func (txInput *TxInput) SetSpend(utxo *UTXO) {
	txInput.SoureId = BytesToHash(utxo.SoureId)
	txInput.SourcePos = utxo.SourcePos
	txInput.Amount = utxo.Amount
	txInput.ScriptPk = utxo.ScriptPk
}

func (txInput *TxInput) ID() (Hash, error) {
	if txInput.ScriptSig == nil {
		return Hash{}, errors.Wrap(inputErr, "txinput ScriptSig empty")
	}

	b := bytes.Join([][]byte{
		txInput.SpendOutputId[:],
		txInput.SoureId[:],
		txInput.ScriptPk,
	}, []byte{})
	return BytesToHash(crypto.Sha256(b)), nil
}

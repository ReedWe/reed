package types

import (
	"bytes"
	"github.com/tybc/crypto"
	"github.com/tybc/errors"
	"github.com/tybc/wallet"
)

var (
	inputErr = errors.New("tx input")
)

type TxInput struct {
	ID Hash `json:"id"`
	Spend
	ScriptSig []byte `json:"scriptSig"` //Sig(Hash(TxInput.Id + tx.Id))
}

func (txInput *TxInput) SetSpend(utxo *UTXO) {
	txInput.SoureId = BytesToHash(utxo.SoureId)
	txInput.SourcePos = utxo.SourcePos
	txInput.Amount = utxo.Amount
	txInput.ScriptPk = utxo.ScriptPk
}

func (txInput *TxInput) SetID() error {
	if txInput.ScriptSig == nil {
		return errors.Wrap(inputErr, "ScriptSig empty")
	}
	b := bytes.Join([][]byte{
		txInput.SpendOutputId[:],
		txInput.SoureId[:],
		txInput.ScriptPk,
	}, []byte{})
	txInput.ID = BytesToHash(crypto.Sha256(b))
	return nil
}

func (txInput *TxInput) SetScriptSig(wt *wallet.Wallet, txId *Hash) error {
	message := bytes.Join([][]byte{
		txInput.ID[:],
		(*txId)[:],
	}, []byte{})

	sig := crypto.Sign(wt.Priv, message)

	//scriptSig = <signature> <public key>
	txInput.ScriptSig = bytes.Join([][]byte{
		sig,
		wt.Pub,
	}, []byte{})
	return nil
}

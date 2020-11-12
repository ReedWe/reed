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

func (txInput *TxInput) GenerateID() (*Hash, error) {
	if txInput.ScriptSig == nil {
		return nil, errors.Wrap(inputErr, "ScriptSig empty")
	}
	b := bytes.Join([][]byte{
		txInput.SpendOutputId.Bytes(),
		txInput.SoureId.Bytes(),
		txInput.ScriptPk,
	}, []byte{})

	h := BytesToHash(crypto.Sha256(b))
	return &h, nil
}

func (txInput *TxInput) GenerateScriptSig(wt *wallet.Wallet, txId *Hash) (*[]byte, error) {
	message := bytes.Join([][]byte{
		txInput.ID.Bytes(),
		(*txId).Bytes(),
	}, []byte{})

	sig := crypto.Sign(wt.Priv, message)

	//scriptSig = <signature> <public key>
	scriptSig := bytes.Join([][]byte{
		sig,
		wt.Pub,
	}, []byte{})
	return &scriptSig, nil
}

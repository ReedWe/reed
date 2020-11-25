package types

import (
	"bytes"
	"encoding/binary"
	"github.com/tybc/crypto"
	"github.com/tybc/errors"
	"github.com/tybc/vm/vmcommon"
	"github.com/tybc/wallet"
)

var (
	inputErr = errors.New("transaction input")
)

type TxInput struct {
	ID        Hash `json:"-"`
	Spend
	ScriptSig []byte `json:"scriptSig"`
}

func (i *TxInput) SetSpend(utxo *UTXO) {
	i.SoureId = utxo.SoureId
	i.SourcePos = utxo.SourcePos
	i.Amount = utxo.Amount
	i.ScriptPk = utxo.ScriptPk
}

func (i *TxInput) GenerateID() Hash {
	split := []byte(":")
	var sourcePosByte = make([]byte, 4)
	binary.LittleEndian.PutUint32(sourcePosByte, i.SourcePos)

	var amountByte = make([]byte, 8)
	binary.LittleEndian.PutUint64(amountByte, i.Amount)

	b := bytes.Join([][]byte{
		i.SpendOutputId.Bytes(),
		split,
		i.SoureId.Bytes(),
		split,
		sourcePosByte,
		split,
		amountByte,
		split,
		i.ScriptPk,
	}, []byte{})

	h := BytesToHash(crypto.Sha256(b))
	return h
}

func (i *TxInput) GenerateScriptSig(wt *wallet.Wallet, txId *Hash) (*[]byte, error) {
	message := bytes.Join([][]byte{
		i.ID.Bytes(),
		(*txId).Bytes(),
	}, []byte{})

	sig := crypto.Sign(wt.Priv, message)

	//scriptSig = <signature> <public key>
	scriptSig := bytes.Join([][]byte{
		{byte(vmcommon.OpPushData64)},
		sig,
		{byte(vmcommon.OpPushData32)},
		wt.Pub,
	}, []byte{})
	return &scriptSig, nil
}

func (i *TxInput) ValidateID() error {
	expect := i.GenerateID()
	if !i.ID.HashEqual(expect) {
		return errors.Wrapf(inputErr, "ID not equal. expect %x. actual %x.", expect, i.ID)
	}
	return nil
}

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

func (ti *TxInput) SetSpend(utxo *UTXO) {
	ti.SoureId = BytesToHash(utxo.SoureId)
	ti.SourcePos = utxo.SourcePos
	ti.Amount = utxo.Amount
	ti.ScriptPk = utxo.ScriptPk
}

func (ti *TxInput) GenerateID() Hash {
	split := []byte(":")
	var sourcePosByte = make([]byte, 4)
	binary.LittleEndian.PutUint32(sourcePosByte, ti.SourcePos)

	var amountByte = make([]byte, 8)
	binary.LittleEndian.PutUint64(amountByte, ti.Amount)

	b := bytes.Join([][]byte{
		ti.SpendOutputId.Bytes(),
		split,
		ti.SoureId.Bytes(),
		split,
		sourcePosByte,
		split,
		amountByte,
		split,
		ti.ScriptPk,
	}, []byte{})

	h := BytesToHash(crypto.Sha256(b))
	return h
}

func (ti *TxInput) GenerateScriptSig(wt *wallet.Wallet, txId *Hash) (*[]byte, error) {
	message := bytes.Join([][]byte{
		ti.ID.Bytes(),
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

func (ti *TxInput) ValidateID() error {
	expect := ti.GenerateID()
	if !ti.ID.HashEqual(expect) {
		return errors.Wrapf(inputErr, "ID not equal. expect %x. actual %x.", expect, ti.ID)
	}
	return nil
}

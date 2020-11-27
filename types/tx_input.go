// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package types

import (
	"bytes"
	"github.com/reed/common/byteutil/byteconv"
	"github.com/reed/crypto"
	"github.com/reed/errors"
	"github.com/reed/vm/vmcommon"
	"github.com/reed/wallet"
)

var (
	inputErr = errors.New("transaction input")
)

type TxInput struct {
	ID Hash `json:"-"`
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
	b := bytes.Join([][]byte{
		i.SpendOutputId.Bytes(),
		split,
		i.SoureId.Bytes(),
		split,
		byteconv.Uint64ToBytes(i.SourcePos),
		split,
		byteconv.Uint64ToBytes(i.Amount),
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

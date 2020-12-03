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
)

type TxOutput struct {
	ID         Hash   `json:"-"`
	IsCoinBase bool   `json:"isCoinBase"`
	Address    []byte `json:"address"`
	Amount     uint64 `json:"amount"`
	ScriptPk   []byte `json:"scriptPK"`
}

var (
	outpuErr = errors.New("transaction output error")
)

func NewTxOutput(isCoinBase bool, address []byte, amount uint64) *TxOutput {
	o := &TxOutput{
		IsCoinBase: isCoinBase,
		Address:    address,
		Amount:     amount,
	}
	o.ScriptPk = o.GenerateLockingScript()
	o.ID = o.GenerateID()
	return o
}

func (o *TxOutput) GenerateID() Hash {
	split := []byte(":")

	data := bytes.Join([][]byte{
		byteconv.BoolToByte(o.IsCoinBase),
		split,
		o.Address,
		split,
		byteconv.Uint64ToByte(o.Amount),
		split,
		o.ScriptPk,
	}, []byte{})

	return BytesToHash(crypto.Sha256(data))
}

func (o *TxOutput) GenerateLockingScript() []byte {
	return vmcommon.BuildP2PKHScript(crypto.Sha256(o.Address))
}

func (o *TxOutput) ValidateID() error {
	expect := o.GenerateID()
	if !o.ID.HashEqual(expect) {
		return errors.Wrapf(outpuErr, "ID not equal. expect %x. actual %x.", expect, o.ID)
	}
	return nil
}

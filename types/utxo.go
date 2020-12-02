// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package types

import (
	"bytes"
	"github.com/reed/common/byteutil/byteconv"
	"github.com/reed/crypto"
)

type UTXO struct {
	ID         Hash   `json:"utxoId"`
	OutputId   Hash   `json:"outputId"`
	SourceId   Hash   `json:"sourceId"`
	IsCoinbase bool   `json:"isCoinbase"`
	SourcePos  uint64 `json:"sourcePos"`
	Amount     uint64 `json:"amount"`
	Address    []byte `json:"address"`
	ScriptPk   []byte `json:"scriptPK"`
}

func NewUtxo(outputId Hash, sourceId Hash, isCoinbase bool, sourcePos uint64, amount uint64, address []byte, scriptPK []byte) *UTXO {
	u := &UTXO{
		OutputId:   outputId,
		SourceId:   sourceId,
		IsCoinbase: isCoinbase,
		SourcePos:  sourcePos,
		Amount:     amount,
		Address:    address,
		ScriptPk:   scriptPK,
	}
	u.ID = u.GenerateID()
	return u
}

func (u *UTXO) GenerateID() Hash {
	split := []byte(":")
	data := bytes.Join([][]byte{
		u.OutputId.Bytes(),
		split,
		u.SourceId.Bytes(),
		split,
		byteconv.BoolToByte(u.IsCoinbase),
		split,
		byteconv.Uint64ToByte(u.SourcePos),
		split,
		byteconv.Uint64ToByte(u.Amount),
		split,
		u.Address,
		split,
		u.ScriptPk,
	}, []byte{})
	return BytesToHash(crypto.Sha256(data))
}

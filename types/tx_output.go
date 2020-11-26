package types

import (
	"bytes"
	"encoding/binary"
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

func (o *TxOutput) GenerateID() Hash {
	split := []byte(":")

	isCoinBaseByte := []byte{1}
	if !o.IsCoinBase {
		isCoinBaseByte = []byte{0}
	}

	var amountByte = make([]byte, 8)
	binary.LittleEndian.PutUint64(amountByte, o.Amount)

	datas := bytes.Join([][]byte{
		isCoinBaseByte,
		split,
		o.Address,
		split,
		amountByte,
		split,
		o.ScriptPk,
	}, []byte{})

	return BytesToHash(crypto.Sha256(datas))
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


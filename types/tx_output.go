package types

import (
	"bytes"
	"encoding/binary"
	"github.com/tybc/crypto"
	"github.com/tybc/errors"
	"github.com/tybc/vm/vmcommon"
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

func (tot *TxOutput) GenerateID() Hash {
	split := []byte(":")

	isCoinBaseByte := []byte{1}
	if !tot.IsCoinBase {
		isCoinBaseByte = []byte{0}
	}

	var amountByte = make([]byte, 8)
	binary.LittleEndian.PutUint64(amountByte, tot.Amount)

	datas := bytes.Join([][]byte{
		isCoinBaseByte,
		split,
		tot.Address,
		split,
		amountByte,
		split,
		tot.ScriptPk,
	}, []byte{})

	return BytesToHash(crypto.Sha256(datas))
}

func (tot *TxOutput) GenerateLockingScript() []byte {
	return vmcommon.BuildP2PKHScript(crypto.Sha256(tot.Address))
}

func (tot *TxOutput) ValidateID() error {
	expect := tot.GenerateID()
	if !tot.ID.HashEqual(expect) {
		return errors.Wrapf(outpuErr, "ID not equal. expect %x. actual %x.", expect, tot.ID)
	}
	return nil
}


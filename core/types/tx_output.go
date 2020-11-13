package types

import (
	"bytes"
	"encoding/binary"
	"github.com/tybc/crypto"
	"github.com/tybc/vm/vmcommon"
)

type TxOutput struct {
	ID         Hash   `json:"id"`
	IsCoinBase bool   `json:"isCoinBase"`
	Address    []byte `json:"address"`
	Amount     uint64 `json:"amount"`
	ScriptPk   []byte `json:"scriptPK"`
}

func (output *TxOutput) GenerateID() Hash {
	//TODO maybe need: len(output.ScriptPK ) >0
	split := []byte(":")

	isCoinBaseByte := []byte{1}
	if !output.IsCoinBase {
		isCoinBaseByte = []byte{0}
	}

	var amountByte = make([]byte, 8)
	binary.LittleEndian.PutUint64(amountByte, output.Amount)

	datas := bytes.Join([][]byte{
		isCoinBaseByte,
		split,
		output.Address,
		split,
		amountByte,
		split,
		output.ScriptPk,
	}, []byte{})

	return BytesToHash(crypto.Sha256(datas))
}

func (output *TxOutput) GenerateLockingScript() []byte {
	return vmcommon.BuildP2PKHScript(crypto.Sha256(output.Address))
}

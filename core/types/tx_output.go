package types

import (
	"bytes"
	"github.com/tybc/crypto"
	"github.com/tybc/vm"
)

// id = Hash(tx.id + address)
type TxOutput struct {
	ID         Hash   `json:"id"`
	IsCoinBase bool   `json:"isCoinBase"`
	Address    []byte `json:"address"`
	Amount     uint64 `json:"amount"`
	ScriptPk   []byte `json:"scriptPK"`
}

func (output *TxOutput) SetID(txId *Hash) {
	b := bytes.Join([][]byte{
		(*txId).Bytes(),
		output.Address,
	}, []byte{})
	output.ID = BytesToHash(b)
}

func (output *TxOutput) SetLockingScript() error {
	output.ScriptPk = vm.BuildP2PKHScript(crypto.Sha256(output.Address))
	return nil
}

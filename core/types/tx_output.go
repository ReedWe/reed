package types

import (
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

func (output *TxOutput) GenerateID(txId *Hash) *Hash {
	h := BytesToHash(crypto.Sha256((*txId).Bytes(), output.Address))
	return &h
}

func (output *TxOutput) GenerateLockingScript() []byte {
	return vm.BuildP2PKHScript(crypto.Sha256(output.Address))
}

package types

import (
	"bytes"
	"github.com/tybc/errors"
)

// id = Hash(tx.id + address)
type TxOutput struct {
	ID         Hash
	IsCoinBase bool
	Address    []byte
	Amount     uint64
	ScriptPk   []byte
}

var (
	outputErr = errors.New("tx output")
)

func (output *TxOutput) SetID(txId *Hash) {
	b := bytes.Join([][]byte{
		(*txId)[:],
		output.Address,
	}, []byte{})
	output.ID = BytesToHash(b)
}

func (output *TxOutput) SetLockingScript() error {

	return nil
}

package types

import (
	_ "github.com/reed/crypto"
)

// id = Hash(field1,field2,...)
type Spend struct {
	SpendOutputId Hash   `json:"spendOutputId"`
	SoureId       Hash   `json:"-"` //rel utxo.SourceId
	SourcePos     uint32 `json:"-"` //rel utxo.SourcePos
	Amount        uint64 `json:"-"` //rel utxo.Amount
	ScriptPk      []byte `json:"-"` //rel utxo.ScriptPK
}

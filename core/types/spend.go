package types

import (
	_ "github.com/tybc/crypto"
)

// id = Hash(field1,field2,...)
type Spend struct {
	SpendOutputId Hash   `json:"spendOutputId"`
	SoureId       Hash   `json:"sourceId"`  //rel utxo.SourceId
	SourcePos     uint32 `json:"sourcePos"` //rel utxo.SourcePos
	Amount        uint64 `json:"amount"`    //rel utxo.Amount
	ScriptPk      []byte `json:"scriptPK"`  //rel utxo.ScriptPK
}

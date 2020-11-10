package types

import (
	_ "github.com/tybc/crypto"
)

// id = Hash(field1,field2,...)
type Spend struct {
	SpendOutputId Hash
	SoureId       Hash   //rel utxo.SourceId
	SourcePos     uint32 //rel utxo.SourcePos
	Amount        uint64 //rel utxo.Amount
	ScriptPk      []byte //rel utxo.ScriptPK
}

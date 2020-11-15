package types

type UTXO struct {
	OutputId  Hash   `json:"outputId"`
	SoureId   Hash   `json:"sourceId"`
	SourcePos uint32 `json:"sourcePos"`
	Amount    uint64 `json:"amount"`
	Address   []byte `json:"address"`
	ScriptPk  []byte `json:"scriptPK"`
}


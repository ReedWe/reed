package types

type TxInput struct {
	TxId          Hash
	Position      uint64
	SpendOutputId Hash

	ScriptSig []byte
}

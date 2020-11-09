package types

type TxInput struct {
	TxId          Hash
	Position      uint32
	SpendOutputId Hash

	ScriptSig []byte
}

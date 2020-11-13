package types


type SubmitTxRequest struct {
	Password  string      `json:"wallet_password"`
	TxInputs  []ReqInput  `json:"tx_inputs"`
	TxOutputs []ReqOutput `json:"tx_outputs"`
}

type ReqInput struct {
	SpendOutputId string `json:"spend_output_id"`
}

type ReqOutput struct {
	Address string `json:"address"`
	Amount  uint64 `json:"amount"`
}

type SumbitTxResponse struct {
	TxId string `json:"tx_id"`
}
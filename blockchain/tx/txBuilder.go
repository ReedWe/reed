package tx

import (
	"encoding/hex"
	"github.com/tybc/blockchain"
	"github.com/tybc/blockchain/validation"
	"github.com/tybc/errors"
)

var (
	ErrSubmitTx = errors.New("sumbit tx")
)

type SubmitTxRequest struct {
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

func SubmitTx(chain *blockchain.Chain, tx *SubmitTxRequest) (*SumbitTxResponse, error) {

	if len(tx.TxInputs) == 0 {
		return nil, errors.WithDetail(ErrSubmitTx, "no input data")
	}

	if len(tx.TxOutputs) == 0 {
		return nil, errors.WithDetail(ErrSubmitTx, "no output data")
	}

	var outputs = make([][]byte, len(tx.TxInputs))
	for i, ti := range tx.TxInputs {
		b, err := hex.DecodeString(ti.SpendOutputId)
		if err != nil {
			return nil, errors.WithDetail(ErrSubmitTx, "invalid spend_output_id format")
		}
		outputs[i] = b
	}

	//TODO check if exist on txpool

	if err := validation.CheckUtxoExists(&chain.Store, &outputs); err != nil {
		return nil, err
	}

	return nil, nil
}

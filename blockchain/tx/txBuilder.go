package tx

import (
	"github.com/tybc/core/types"
	"github.com/tybc/blockchain/validation"
)

//type BuildResult struct {
//	Transaction TxResult `json:"transaction"`
//}

type BuildRequest struct {
	TxInputs  []types.TxInput
	TxOutputs []types.TxOutput
}

type ReqInput struct {
	SpendTxId     types.Hash
	SpendPosition uint64
	Signature     string
}

type ReqOutput struct {
	Address types.Hash
	Amount  uint64
}

type BuildResponse struct {
	TxId types.Hash
}

func Build(tx *BuildRequest) (*BuildResponse, error) {

	if len(tx.TxInputs) == 0 || len(tx.TxOutputs) == 0 {
		return nil, nil
	}

	err := validation.CheckUtxoExists(&tx.TxOutputs)

	if err != nil {
		return nil, nil
	}

	return nil, nil
}

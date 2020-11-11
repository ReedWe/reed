package txbuilder

import (
	"encoding/hex"
	"github.com/tybc/blockchain"
	"github.com/tybc/core"
	"github.com/tybc/core/types"
	"github.com/tybc/errors"
	"github.com/tybc/log"
	"github.com/tybc/wallet"
)

var (
	ErrSubmitTx = errors.New("sumbit tx")
)

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

func (req *SubmitTxRequest) MapTx() (*types.Tx, error) {

	inputs := make([]types.TxInput, len(req.TxInputs))
	outputs := make([]types.TxOutput, len(req.TxOutputs))

	for _, inp := range req.TxInputs {
		b, err := hex.DecodeString(inp.SpendOutputId)
		if err != nil {
			return nil, errors.WithDetail(ErrSubmitTx, "invalid spend_output_id format")
		}
		inputs = append(inputs, types.TxInput{
			Spend: types.Spend{SpendOutputId: types.BytesToHash(b)},
		})
	}

	for _, iop := range req.TxOutputs {
		addr, err := hex.DecodeString(iop.Address)
		if err != nil {
			return nil, errors.WithDetail(ErrSubmitTx, "invalid output.address format")
		}

		outputs = append(outputs, types.TxOutput{
			IsCoinBase: false,
			Address:    addr,
			Amount:     iop.Amount,
		})
	}

	tx := &types.Tx{
		TxInput:  inputs,
		TxOutput: outputs,
	}

	return tx, nil
}

func SubmitTx(chain *blockchain.Chain, reqTx *SubmitTxRequest) (*SumbitTxResponse, error) {

	if len(reqTx.TxInputs) == 0 {
		return nil, errors.WithDetail(ErrSubmitTx, "no input data")
	}

	if len(reqTx.TxOutputs) == 0 {
		return nil, errors.WithDetail(ErrSubmitTx, "no output data")
	}

	// get wallet
	wt, err := wallet.My(reqTx.Password)
	if err != nil {
		log.Logger.Infof("pub %s", wt.Pub)
		return nil, err
	}

	// request data map to tx
	tx, err := reqTx.MapTx()
	if err != nil {
		return nil, err
	}

	//check input rel utxo
	//set spend data
	for _, input := range tx.TxInput {
		if utxo, err := core.GetUtxoByOutputId(&chain.Store, input.SpendOutputId); err != nil {
			return nil, err
		} else {
			input.SetSpend(utxo)
			input.SetID()
		}
	}

	tx.SetID()

	//sign scriptSig
	for _, input := range tx.TxInput {
		input.SetScriptSig(wt, tx.ID)
	}

	//set outpu id
	//locking script
	for _, output := range tx.TxOutput {
		output.SetID(&tx.ID)
	}

	//TODO check if exist on txpool
	return nil, nil
}

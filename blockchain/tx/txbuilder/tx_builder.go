package txbuilder

import (
	"encoding/hex"
	"github.com/tybc/blockchain"
	"github.com/tybc/core/types"
	"github.com/tybc/errors"
	"github.com/tybc/log"
	"github.com/tybc/wallet"
)

var (
	ErrSubmitTx = errors.New("sumbit tx")
)

func MapTx(req *types.SubmitTxRequest) (*types.Tx, error) {

	var inputs []types.TxInput
	var outputs []types.TxOutput

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

func SubmitTx(chain *blockchain.Chain, reqTx *types.SubmitTxRequest) (*types.SumbitTxResponse, error) {

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
	tx, err := MapTx(reqTx)
	if err != nil {
		return nil, err
	}

	//check input rel utxo
	//set spend data
	for _, input := range tx.TxInput {
		if utxo, err := blockchain.GetUtxoByOutputId(&chain.Store, input.SpendOutputId); err != nil {
			return nil, err
		} else {
			input.SetSpend(utxo)
		}
	}

	//tx ID
	txId, err := tx.GenerateID()
	if err != nil {
		return nil, err
	}
	tx.ID = *txId

	for _, input := range tx.TxInput {
		//ScriptSig
		scriptSig, err := input.GenerateScriptSig(wt, &tx.ID)
		if err != nil {
			return nil, err
		}
		input.ScriptSig = *scriptSig

		//ID
		id, err := input.GenerateID()
		if err != nil {
			return nil, err
		}
		input.ID = *id
	}

	//set output id
	//locking script
	for _, output := range tx.TxOutput {
		output.ID = *output.GenerateID(&tx.ID)
		output.ScriptPk = output.GenerateLockingScript()
	}

	//TODO check if exist on txpool
	return nil, nil
}

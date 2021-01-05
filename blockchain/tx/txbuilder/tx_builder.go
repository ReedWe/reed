// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package txbuilder

import (
	"crypto/ed25519"
	"encoding/hex"
	"github.com/reed/blockchain"
	"github.com/reed/blockchain/tx/txpusher"
	"github.com/reed/common/math"
	"github.com/reed/errors"
	"github.com/reed/types"
	"github.com/reed/wallet"
)

var (
	submitTxErr   = errors.New("submit transaction")
	txAssetAmtErr = errors.New("transaction asset amount error")
)

func SubmitTx(chain *blockchain.Chain, reqTx *types.SubmitTxRequest) (*types.SumbitTxResponse, error) {

	if len(reqTx.TxInputs) == 0 {
		return nil, errors.WithDetail(submitTxErr, "no input data")
	}

	if len(reqTx.TxOutputs) == 0 {
		return nil, errors.WithDetail(submitTxErr, "no output data")
	}

	// get wallet
	wt, err := wallet.My(reqTx.Password)
	if err != nil {
		return nil, err
	}

	// request data map to tx
	tx, err := mapTx(reqTx)
	if err != nil {
		return nil, err
	}

	if err = inspectionTx(tx, wt.Pub); err != nil {
		return nil, err
	}

	getUtxo := func(spendOutputId types.Hash) (*types.UTXO, error) {
		return blockchain.GetUtxoByOutputId(chain.Store, spendOutputId)
	}
	if err = tx.Completion(getUtxo); err != nil {
		return nil, err
	}

	txId := tx.GetID()
	// ScriptSig
	for _, input := range tx.TxInput {
		scriptSig, err := input.GenerateScriptSig(wt, &txId)
		if err != nil {
			return nil, err
		}
		input.ScriptSig = *scriptSig
	}

	if err = txpusher.MaybePush(chain, tx); err != nil {
		return nil, err
	}
	return nil, nil
}

func mapTx(req *types.SubmitTxRequest) (*types.Tx, error) {

	var inputs []*types.TxInput
	var outputs []*types.TxOutput

	for _, inp := range req.TxInputs {
		b, err := hex.DecodeString(inp.SpendOutputId)
		if err != nil {
			return nil, errors.WithDetail(submitTxErr, "invalid spend_output_id format")
		}
		inputs = append(inputs, &types.TxInput{
			Spend: types.Spend{SpendOutputId: types.BytesToHash(b)},
		})
	}

	for _, iop := range req.TxOutputs {
		addr, err := hex.DecodeString(iop.Address)
		if err != nil {
			return nil, errors.WithDetail(submitTxErr, "invalid output.address format")
		}

		outputs = append(outputs, types.NewTxOutput(false, addr, iop.Amount))
	}

	tx := &types.Tx{
		TxInput:  inputs,
		TxOutput: outputs,
	}

	return tx, nil
}

func inspectionTx(tx *types.Tx, pub ed25519.PublicKey) error {
	if err := maybeFillSelfOutput(tx, pub); err != nil {
		return err
	}

	newOutputs, err := mergeSameAddrOutput(tx.TxOutput)
	if err != nil {
		return err
	}
	tx.TxOutput = newOutputs

	return nil
}

func maybeFillSelfOutput(tx *types.Tx, pub ed25519.PublicKey) error {
	sumInput, sumOutput, err := tx.IsAssetAmtEqual()
	if err != nil {
		return err
	}
	if sumOutput > sumInput {
		return errors.Wrap(txAssetAmtErr, "not enough inputs amount")
	}
	if sumInput > sumOutput {
		// auto generate self output
		tx.TxOutput = append(tx.TxOutput, types.NewTxOutput(false, pub, sumInput-sumOutput))
	}
	return nil
}

func mergeSameAddrOutput(outputs []*types.TxOutput) ([]*types.TxOutput, error) {
	var newOutputs []*types.TxOutput
	addrMap := map[string]*types.TxOutput{}
	var err error
	for _, output := range outputs {
		if item, ok := addrMap[string(output.Address)]; ok {
			// exist same address output,merge
			item.Amount, err = math.AddUint64(item.Amount, output.Amount)
			if err != nil {
				return nil, err
			}
		} else {
			newOutputs = append(newOutputs, output)
			addrMap[string(output.Address)] = output
		}
	}
	return newOutputs, nil
}

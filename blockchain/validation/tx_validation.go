package validation

import (
	"bytes"
	"github.com/tybc/crypto"
	"github.com/tybc/errors"
	"github.com/tybc/log"
	"github.com/tybc/types"
	"github.com/tybc/vm"
)

var (
	validationInputErr  = errors.New("validation tx:invalid input")
	validationOutputErr = errors.New("validation tx:invalid output")
)

func ValidateTx(tx *types.Tx) error {
	log.Logger.Infof("ValidateTx %v", *tx)

	if err := validateInput(&tx.TxInput, tx.ID); err != nil {
		return err
	}

	if err := validateOutput(&tx.TxOutput); err != nil {
		return err
	}

	if _, _, err := tx.IsAssetAmtEqual(); err != nil {
		return err
	}

	return nil
}

func validateInput(inputs *[]types.TxInput, txId types.Hash) error {
	if len(*inputs) == 0 {
		return errors.Wrapf(validationInputErr, "input data empty")
	}

	spendOutputMap := map[string]*types.TxInput{}
	for _, input := range *inputs {

		key := string(input.SoureId.Bytes()) + string(input.SourcePos)
		if _, ok := spendOutputMap[key]; ok {
			return errors.Wrapf(validationInputErr, "repeat utxos,source id=%x", input.SoureId)
		}

		// notes: utxo check in transaction.Completion

		// verify sign
		signFunc := func(pk []byte, sig []byte) bool {
			message := bytes.Join([][]byte{
				input.ID.Bytes(),
				txId.Bytes(),
			}, []byte{})
			return crypto.Verify(pk, message, sig)
		}

		virtualMachine := vm.NewVirtualMachine(input.ScriptSig, input.ScriptPk, signFunc)
		if err := virtualMachine.Run(); err != nil {
			return err
		}
	}
	return nil
}

func validateOutput(outputs *[]types.TxOutput) error {
	if len(*outputs) == 0 {
		return errors.Wrapf(validationOutputErr, "output data empty")
	}

	for _, output := range *outputs {
		if output.IsCoinBase {
			return errors.Wrapf(validationOutputErr, "not coinbase output %x", output.ID)
		}
	}
	return nil
}

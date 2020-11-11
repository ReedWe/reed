package validation

import (
	"bytes"
	"github.com/tybc/blockchain"
	"github.com/tybc/core/types"
	"github.com/tybc/crypto"
	"github.com/tybc/errors"
	"github.com/tybc/log"
	"github.com/tybc/vm"
)

var (
	validationInputErr  = errors.New("validation tx:invalid input")
	validationOutputErr = errors.New("validation tx:invalid output")
)

func ValidateTx(chain *blockchain.Chain, tx *types.Tx) error {
	log.Logger.Infof("ValidateTx %v", *tx)

	//TODO check exist in tx pool

	if err := validateInput(chain, &tx.TxInput); err != nil {
		return err
	}

	if err := validateOutput(&tx.TxOutput); err != nil {
		return err
	}

	return nil
}

func validateInput(chain *blockchain.Chain, inputs *[]types.TxInput) error {
	for _, input := range *inputs {
		if _, err := blockchain.GetUtxoByOutputId(&chain.Store, input.SpendOutputId); err != nil {
			return errors.Wrap(validationInputErr, err)
		}
		if len(input.ScriptSig) != (64 + 32) {
			return errors.Wrapf(validationInputErr, "invalid scriptSig len(%d).input ID %x", len(input.ScriptSig), input.ID)
		}
		//TODO RUN VM
	}
	return nil
}

func validateOutput(outputs *[]types.TxOutput) error {
	for _, output := range *outputs {
		if len(output.Address) != 32 {
			return errors.Wrapf(validationOutputErr, "invalid output address. len(%d).expect 36", len(output.Address))
		}

		script := output.ScriptPk
		if len(script) != (32 + 4) {
			return errors.Wrapf(validationOutputErr, "invalid ScriptPk len(%d).expect 36", len(script))
		}
		if script[0] != vm.OpDup {
			return errors.Wrapf(validationOutputErr, "ScriptPk first part is not OP_DUP.output ID %x", output.ID)
		}
		if script[1] != vm.OpHash256 {
			return errors.Wrapf(validationOutputErr, "ScriptPk second part is not OP_HASH256.output ID %x", output.ID)
		}

		pubHash := crypto.Sha256(output.Address)
		if !bytes.Equal(script[2:34], pubHash) {
			return errors.Wrapf(validationOutputErr, "ScriptPk third part is not Hash data.output ID %x", output.ID)
		}

		if script[34] != vm.OpEqualverify {
			return errors.Wrapf(validationOutputErr, "ScriptPk fourth part is not OP_EQUALVERIFY.output ID %x", output.ID)
		}
		if script[35] != vm.OpChecksig {
			return errors.Wrapf(validationOutputErr, "ScriptPk fifth part is not OP_CHECKSIG.output ID %x", output.ID)
		}
	}
	return nil
}

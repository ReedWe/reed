package validation

import (
	"bytes"
	"github.com/tybc/blockchain"
	"github.com/tybc/core/types"
	"github.com/tybc/crypto"
	"github.com/tybc/errors"
	"github.com/tybc/log"
	"github.com/tybc/vm/vmcommon"
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

	if _, _, err := tx.IsAssetAmtEqual(); err != nil {
		return err
	}

	return nil
}

func validateInput(chain *blockchain.Chain, inputs *[]types.TxInput) error {
	spendOutputMap := map[string]*types.TxInput{}
	for _, input := range *inputs {
		key := string(input.SoureId.Bytes()) + string(input.SourcePos)
		if _, ok := spendOutputMap[key]; ok {
			return errors.Wrapf(validationInputErr, "repeat utxos,source id=%x", input.SoureId)
		}

		if _, err := blockchain.GetUtxoByOutputId(&chain.Store, input.SpendOutputId); err != nil {
			return errors.Wrap(validationInputErr, err)
		}
		if len(input.ScriptSig) != (64 + 32 + 2) {
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
		if len(script) != (32 + 5) {
			return errors.Wrapf(validationOutputErr, "invalid ScriptPk len(%d).expect 36", len(script))
		}
		if script[0] != vmcommon.OpDup {
			return errors.Wrapf(validationOutputErr, "ScriptPk first part is not OP_DUP.output ID %x", output.ID)
		}
		if script[1] != vmcommon.OpHash256 {
			return errors.Wrapf(validationOutputErr, "ScriptPk second part is not OP_HASH256.output ID %x", output.ID)
		}

		if script[2] != vmcommon.OpPushData32 {
			return errors.Wrapf(validationOutputErr, "ScriptPk third part is not OP_PUSHDATA32.output ID %x", output.ID)
		}

		pubHash := crypto.Sha256(output.Address)
		if !bytes.Equal(script[3:35], pubHash) {
			return errors.Wrapf(validationOutputErr, "ScriptPk fourth part is not Hash data.output ID %x", output.ID)
		}

		if script[35] != vmcommon.OpEqualVerify {
			return errors.Wrapf(validationOutputErr, "ScriptPk fifth part is not OP_EQUALVERIFY.output ID %x", output.ID)
		}
		if script[36] != vmcommon.OpCheckSig {
			return errors.Wrapf(validationOutputErr, "ScriptPk sixth part is not OP_CHECKSIG.output ID %x", output.ID)
		}
	}
	return nil
}

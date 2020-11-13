package txbuilder

import (
	"github.com/tybc/blockchain"
	"github.com/tybc/blockchain/validation"
	"github.com/tybc/types"
)

func MaybePush(chain *blockchain.Chain, tx *types.Tx) error {
	if err := validation.ValidateTx(chain, tx); err != nil {
		return err
	}
	//TODO push into tx pool

	return nil
}

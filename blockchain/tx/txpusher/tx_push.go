package txpusher

import (
	"github.com/tybc/blockchain"
	"github.com/tybc/blockchain/validation"
	"github.com/tybc/log"
	"github.com/tybc/types"
)

func MaybePush(chain *blockchain.Chain, tx *types.Tx) error {
	log.Logger.Infof("receive a new transaction %x", tx.ID)

	existTx, err := chain.Txpool.GetTx(&tx.ID)
	if err != nil {
		return err
	}
	if existTx != nil {
		log.Logger.Infof("transaction exists already (id=%x)", tx.ID)
		return nil
	}

	if err := validation.ValidateTx(chain, tx); err != nil {
		return err
	}
	//TODO push into tx pool

	return nil
}

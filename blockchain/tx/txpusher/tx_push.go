package txpusher

import (
	"github.com/tybc/blockchain"
	"github.com/tybc/blockchain/validation"
	"github.com/tybc/errors"
	"github.com/tybc/log"
	"github.com/tybc/types"
)

var (
	recvTxErr = errors.New("receive transaction error")
)

// receive local transaction and remote transaction
func MaybePush(chain *blockchain.Chain, tx *types.Tx) error {
	log.Logger.Infof("receive a new transaction ID=%x", tx.ID)

	if err := tx.Completion(&chain.Store); err != nil {
		return err
	}

	//tx ID
	txId, err := tx.GenerateID()
	if !txId.HashEqual(tx.ID) {
		return errors.Wrapf(recvTxErr, "txId errors. local=%x remote=%x", txId, tx.ID)
	}

	existTx, err := chain.Txpool.GetTx(&tx.ID)
	if err != nil {
		return err
	}

	if existTx != nil {
		log.Logger.Infof("transaction exists already (id=%x)", tx.ID)
		return nil
	}

	if err := validation.ValidateTx(tx); err != nil {
		return err
	}
	//TODO push into tx pool

	return nil
}

package blockchain

import (
	bm "github.com/reed/blockchain/blockmanager"
	"github.com/reed/blockchain/config"
	"github.com/reed/blockchain/store"
	"github.com/reed/blockchain/txpool"
	"github.com/reed/blockchain/validation"
	"github.com/reed/errors"
	"github.com/reed/log"
	"github.com/reed/types"
	"github.com/sirupsen/logrus"
)

type Chain struct {
	Store            *store.Store
	Txpool           *txpool.Txpool
	BlockManager     *bm.BlockManager
	blockReceptionCh chan *types.RecvWrap
	breakWorkCh      chan struct{}
	isOpen           bool
}

var (
	openChainErr         = errors.New("chain:open chain error")
	processBlockChainErr = errors.New("chain:process block error")
)

func NewChain(s *store.Store) (*Chain, error) {
	tp := txpool.NewTxpool(s)
	highestBlock, err := (*s).GetHighestBlock()
	if err != nil {
		return nil, err
	}

	recvCh := make(chan *types.RecvWrap, 100)
	bwCh := make(chan struct{})

	blockMgr, err := bm.NewBlockManager(s, highestBlock, recvCh)
	if err != nil {
		return nil, err
	}

	c := &Chain{
		Store:            s,
		Txpool:           tp,
		BlockManager:     blockMgr,
		blockReceptionCh: recvCh,
		breakWorkCh:      bwCh,
	}
	return c, nil
}

func (c *Chain) Open() error {
	if c.isOpen {
		return errors.Wrap(openChainErr, "There is already an open channel.")
	}
	go receiveBlock(c, c.blockReceptionCh, c.breakWorkCh)
	log.Logger.Info("★★Chain is open")
	return nil
}

func (c *Chain) Close() {
	c.isOpen = false
	close(c.blockReceptionCh)
	log.Logger.Info("★★Chain is close")

}

func (c *Chain) GetReadBreakWorkChan() <-chan struct{} {
	return c.breakWorkCh
}

func (c *Chain) GetWriteReceptionChan() chan<- *types.RecvWrap {
	return c.blockReceptionCh
}

// save new block and broadcast
func (c *Chain) ProcessNewBlock(block *types.Block) error {
	if err := validation.ValidateBlockHeader(block, c.BlockManager.HighestBlock(), c.BlockManager); err != nil {
		return err
	}
	exists, err := c.BlockManager.AddNewBlock(block)
	if exists {
		return nil
	}
	if err != nil {
		return err
	}

	// TODO broadcast new blockmanager

	if err := ProcessUtxoForSaveBlock(c.Store, block); err != nil {
		return errors.Wrapf(processBlockChainErr, err.Error())
	}

	// remove processed transaction in pool
	c.Txpool.RemoveTransactions(block.Transactions)

	return nil
}

func receiveBlock(chain *Chain, receptionCh <-chan *types.RecvWrap, stopWorkCh chan<- struct{}) {
	for item := range receptionCh {
		block := item.Block
		log.Logger.WithFields(logrus.Fields{"height": block.Height, "hash": block.GetHash().ToString(), "SendBreakWork": item.SendBreakWork}).Info("receive a new block")
		if err := chain.ProcessNewBlock(block); err != nil {
			log.Logger.WithField("blockHash", block.GetHash().ToString()).Error(err)
		} else {
			if item.SendBreakWork && config.Default.Mining {
				stopWorkCh <- struct{}{}
			}
		}
	}
	log.Logger.Info("receiveBlock is stop.")
}

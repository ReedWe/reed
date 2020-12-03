// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package miner

import (
	"bytes"
	"fmt"
	bc "github.com/reed/blockchain"
	bm "github.com/reed/blockchain/blockmanager"
	"github.com/reed/consensus/pow"
	"github.com/reed/errors"
	"github.com/reed/log"
	"github.com/reed/types"
	"github.com/reed/wallet"
	"strconv"
	"sync"
)

var (
	startErr = errors.New("miner failed to start")
	workErr  = errors.New("miner failed to work")
)

const (
	maxTries = ^uint64(0)
)

type Miner struct {
	sync.Mutex
	wallet           *wallet.Wallet
	chain            *bc.Chain
	working          bool
	blockManager     *bm.BlockManager
	blockReceptionCh <-chan *types.Block
	blockSendCh      chan<- *types.Block
	stopWorkCh       <-chan struct{}
}

func NewMiner(c *bc.Chain, w *wallet.Wallet, sch <-chan *types.Block) *Miner {
	return &Miner{
		wallet:           w,
		chain:            c,
		working:          false,
		blockReceptionCh: sch,
	}
}

func (m *Miner) Start() error {
	m.Lock()
	defer m.Unlock()

	if m.working {
		return errors.Wrap(startErr, "miner has started.")
	}

	m.work()

	return nil
}

func (m *Miner) work() {
	highest, err := m.chain.BlockManager.HighestBlock()
	if err != nil {
		log.Logger.Error(workErr, err)
		return
	}
	block, err := m.buildBlock(highest)
	if err != nil {
		log.Logger.Error(workErr, err)
		return
	}

	for {
		born, stop := m.generateBlock(block)
		if born {
			fmt.Printf("mint a new block %x %v \n", block.GetHash(), block)
			break
			m.blockManager.AddNewBlock(block)
			//broadcast new blockmanager

			newBlock, err := m.buildBlock(highest)
			if err != nil {
				log.Logger.Error(workErr, err)
				break
			}
			block = newBlock
		}
		if stop {
			break
		}
	}
}

func (m *Miner) generateBlock(block *types.Block) (born bool, stop bool) {
	extraNonce := uint64(0)
loop:
	for {
		select {
		case rblock := <-m.blockReceptionCh:
			log.Logger.Infof("Received a blockmanager from blockReception channel.id=%x", rblock.GetHash())
			// receive a new block from remote node
			// block = fetch laest block
			block = rblock
			born = true
			break loop
		case <-m.stopWorkCh:
			log.Logger.Info("Received a stop single,stop miner...")
			stop = true
			break loop
		default:
			//just for no block,do nothing
		}

		if pow.CheckProofOfWork(block.BigNumber, block.GetHash()) {
			born = true
			break loop
		} else {
			if block.Nonce == maxTries {
				//reset nonce
				block.Nonce = 0

				//change coinbase tx's scriptSig and continue
				extraNonce++
				m.incrementExtraNonce(extraNonce, block)
			} else {
				block.Nonce++
			}
		}
	}
	return
}

func (m *Miner) buildBlock(pre *types.Block) (*types.Block, error) {
	var newBlock *types.Block
	if pre == nil {
		newBlock = types.GenesisBlock()
		newBlock.BigNumber = pow.DifficultyLimit()
	} else {
		newBlock = &types.Block{
			BlockHeader:  *pre.Copy(),
			Transactions: []*types.Tx{},
		}
	}

	txs := m.chain.Txpool.GetTxs()
	cbTx, err := types.NewCoinbaseTx(newBlock.Height, m.wallet.Pub, bc.CalcCoinbaseAmt(newBlock.Height))
	if err != nil {
		return nil, err
	}
	if len(txs) == 0 {
		txs = append(txs, cbTx)
	} else {
		txs = append(txs, nil)
		copy(txs[1:], txs[:len(txs)-1])
		txs[0] = cbTx
	}
	newBlock.Transactions = txs

	//recalculate difficulty
	newBlock.BigNumber = pow.GetNextDifficulty(newBlock, m.blockManager.GetAncestor)
	return newBlock, nil
}

func (m *Miner) incrementExtraNonce(extraNonce uint64, cblock *types.Block) {
	cblock.Transactions[0].TxInput[0].ScriptSig = bytes.Join([][]byte{cblock.Transactions[0].TxInput[0].ScriptSig, []byte(strconv.FormatUint(extraNonce, 10))}, []byte{})
}

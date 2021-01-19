// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package miner

import (
	"bytes"
	"fmt"
	bc "github.com/reed/blockchain"
	"github.com/reed/blockchain/merkle"
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
	blockReceptionCh chan<- *types.RecvWrap
	breakWorkCh      <-chan struct{}
	stopCh           chan struct{}
}

func NewMiner(c *bc.Chain, w *wallet.Wallet, rCh chan<- *types.RecvWrap, bCh <-chan struct{}) *Miner {
	return &Miner{
		wallet:           w,
		chain:            c,
		working:          false,
		blockReceptionCh: rCh,
		breakWorkCh:      bCh,
	}
}

func (m *Miner) Start() error {
	m.Lock()
	defer m.Unlock()

	if m.working {
		return errors.Wrap(startErr, "Miner has started.")
	}
	m.working = true
	go m.work()
	fmt.Println("★★Miner Server Start")
	return nil
}

func (m *Miner) Stop() {
	m.Lock()
	defer m.Unlock()

	m.working = false
	fmt.Println("★★Miner Server Stop")
}

func (m *Miner) fetchBlock() (*types.Block, error) {
	highest := m.chain.BlockManager.HighestBlock()
	block, err := m.buildBlock(highest)
	if err != nil {
		log.Logger.Error(workErr, err)
		return nil, err
	}
	return block, nil
}

func (m *Miner) work() {
	block, err := m.fetchBlock()
	if err != nil {
		log.Logger.Fatal(workErr, err)
		return
	}
	for {
		select {
		case <-m.stopCh:
			log.Logger.Info("mining work is stop")
		default:
		}

		repack := m.generateBlock(block)
		if repack {
			block, err = m.fetchBlock()
			if err != nil {
				log.Logger.Fatal(workErr, err)
				break
			}
			log.Logger.Info("receive from remote or reorganize chan,repack block complete.")
		} else {
			block, err = m.buildBlock(block)
			if err != nil {
				log.Logger.Error(workErr, err)
				break
			}
			log.Logger.Info("mined a block,rebuild a new block complete.")
		}
	}
}

func (m *Miner) generateBlock(block *types.Block) (repack bool) {
	extraNonce := uint64(0)
loop:
	for {
		select {
		case <-m.breakWorkCh:
			log.Logger.Info("Received a break single,stop mining.")
			repack = true
			break loop
		default:
			// just for no block,do nothing
		}

		if pow.CheckProofOfWork(block.BigNumber, block.GetHash()) {
			m.blockReceptionCh <- &types.RecvWrap{SendBreakWork: false, Block: block}
			break loop
		} else {
			if block.Nonce == maxTries {
				// reset nonce
				block.Nonce = 0

				// change coinbase tx's scriptSig and continue
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
		newBlock = types.GetGenesisBlock()
	} else {
		newBlock = &types.Block{
			BlockHeader:  *pre.Copy(),
			Transactions: []*types.Tx{},
		}
	}
	newBlock.BigNumber = pow.DifficultyLimit()

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

	// recalculate difficulty
	newBlock.BigNumber = pow.GetDifficulty(newBlock, m.chain.BlockManager.GetAncestor)
	// set tx merkle root
	newBlock.MerkleRootHash = merkle.ComputeMerkleRoot(newBlock.Transactions)
	return newBlock, nil
}

// when the nonce reaches the maximum value,change the scriptSig value of coinbase transaction
// and reset:nonce=0
func (m *Miner) incrementExtraNonce(extraNonce uint64, b *types.Block) {
	b.Transactions[0].TxInput[0].ScriptSig = bytes.Join([][]byte{b.Transactions[0].TxInput[0].ScriptSig, []byte(strconv.FormatUint(extraNonce, 10))}, []byte{})
	// recompute merkle root
	b.MerkleRootHash = merkle.ComputeMerkleRoot(b.Transactions)
}

// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package miner

import (
	"bytes"
	"fmt"
	"github.com/reed/blockchain/block"
	"github.com/reed/consensus/pow"
	"github.com/reed/errors"
	"github.com/reed/log"
	"strconv"
	"sync"
)

var (
	startErr = errors.New("miner start error")
)

const (
	maxTries = ^uint64(0)
)

type Miner struct {
	sync.Mutex
	working          bool
	blockReceptionCh <-chan *block.Block
	blockSendCh      chan<- *block.Block
	stopWorkCh       <-chan struct{}
}

func NewMiner(submitCh <-chan *block.Block) *Miner {
	return &Miner{
		working:          false,
		blockReceptionCh: submitCh,
	}
}

func (m *Miner) Start() error {
	m.Lock()
	defer m.Unlock()

	if m.working {
		return errors.Wrap(startErr, "miner has started.")
	}

	go m.work()

	return nil
}

func (m *Miner) work() {
	var block block.Block
	//calc difficulty
	block.BigNumber = pow.GetNextDifficulty(&block)

	extraNonce := uint64(0)

loop:
	for {
		select {
		case b := <-m.blockReceptionCh:
			// receive a new block from remote node
			// block = fetch laest block
			fmt.Println(b)
		case <-m.stopWorkCh:
			log.Logger.Info("receive a stop single,stop miner...")
			break loop
		default:
			//just for no block,do nothing
		}

		if pow.CheckProofOfWork(block.BigNumber, block.GetHash()) {
			//broadcast new block
		} else {
			if block.Nonce == maxTries {
				//reset nonce
				block.Nonce = 0

				//change coinbase tx's scriptSig and continue
				extraNonce++
				m.incrementExtraNonce(extraNonce, &block)
			} else {
				block.Nonce++
			}
		}
	}
}

func (m *Miner) incrementExtraNonce(extraNonce uint64, cblock *block.Block) {
	txs := *cblock.Transactions
	txs[0].TxInput[0].ScriptSig = bytes.Join([][]byte{txs[0].TxInput[0].ScriptSig, []byte(strconv.FormatUint(extraNonce, 10))}, []byte{})
}

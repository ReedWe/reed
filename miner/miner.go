// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package miner

import (
	"bytes"
	"fmt"
	"github.com/reed/consensus/pow"
	"github.com/reed/errors"
	"github.com/reed/log"
	"github.com/reed/types"
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
	blockReceptionCh <-chan *types.Block
	blockSendCh      chan<- *types.Block
	stopWorkCh       <-chan struct{}
}

func NewMiner(submitCh <-chan *types.Block) *Miner {
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
	var cblock types.Block

	extraNonce := uint64(0)

loop:
	for {
		select {
		case b := <-m.blockReceptionCh:
			// receive a new block from remote node
			// cblock = fetch laest block
			fmt.Println(b)
		case <-m.stopWorkCh:
			log.Logger.Info("receive a stop single,stop miner...")
			break loop
		default:
			//just for no block,do nothing
		}

		if pow.CheckProofOfWork(cblock.Bits, cblock.GetHash()) {
			//broadcast new block
		} else {
			if cblock.Nonce == maxTries {
				//reset nonce
				cblock.Nonce = 0

				//change coinbase tx's scriptSig and continue
				extraNonce++
				m.incrementExtraNonce(extraNonce, &cblock)
			} else {
				cblock.Nonce++
			}
		}
	}
}

func (m *Miner) incrementExtraNonce(extraNonce uint64, cblock *types.Block) {
	txs := *cblock.Transactions
	txs[0].TxInput[0].ScriptSig = bytes.Join([][]byte{txs[0].TxInput[0].ScriptSig, []byte(strconv.FormatUint(extraNonce, 10))}, []byte{})
}

// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package miner

import (
	"bytes"
	"fmt"
	"github.com/reed/consensus/pow"
	"github.com/reed/errors"
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
	working  bool
	submitCh <-chan *types.Block
	winCh    chan<- *types.Block
}

func NewMiner(submitCh <-chan *types.Block) *Miner {
	return &Miner{
		working:  false,
		submitCh: submitCh,
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

	//获取最新的区块信息，赋值给cblock

	tries := maxTries
	extraNonce := uint64(0)

	for {
		select {
		case b := <-m.submitCh:
			fmt.Println(b)

		default:
			//just for no block,do nothing
		}

		if pow.CheckProofOfWork(cblock.Bits, cblock.GetHash()) {
			//挖矿成功广播区块
		} else {
			if tries == 0 {
				//reset
				tries = maxTries
				cblock.Nonce = 0

				extraNonce++
				m.incrementExtraNonce(extraNonce, &cblock)
			} else {
				cblock.Nonce++
				tries--
			}
		}
	}
}

func (m *Miner) incrementExtraNonce(extraNonce uint64, cblock *types.Block) {
	txs := *cblock.Transactions
	txs[0].TxInput[0].ScriptSig = bytes.Join([][]byte{txs[0].TxInput[0].ScriptSig, []byte(strconv.FormatUint(extraNonce, 10))}, []byte{})
}

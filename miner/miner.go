// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package miner

import (
	"fmt"
	"github.com/reed/consensus/pow"
	"github.com/reed/errors"
	"github.com/reed/types"
	"math/big"
	"sync"
)

var (
	startErr = errors.New("miner start error")
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
	var cblock types.BlockHeader

	//获取最新的区块信息，赋值给cblock

	for {
		select {
		case b := <-m.submitCh:
			fmt.Println(b)
			cblock = b.BlockHeader
		default:
			//just for no block,do nothing
		}

		if pow.CheckProofOfWork(cblock.Bits, cblock.GetHash()) {
			//挖矿成功广播区块
		} else {
			//判断 nonce是否已达到最大值
			cblock.Nonce.Add(&cblock.Nonce, big.NewInt(1))
		}
	}

}

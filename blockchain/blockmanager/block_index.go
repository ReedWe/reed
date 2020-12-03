// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package blockmanager

import (
	"github.com/reed/blockchain/store"
	"github.com/reed/log"
	"github.com/reed/types"
)

type BlockIndex struct {
	store *store.Store
	index map[types.Hash]*types.Block
	main  []*types.Block
}

const (
	mainArrayInterval = 24 * 6
)

func NewBlockIndex(s *store.Store, highestBlock *types.Block) (*BlockIndex, error) {
	h := uint64(0)
	c := uint64(mainArrayInterval)

	if highestBlock != nil {
		h = highestBlock.Height
	}
	if h > mainArrayInterval {
		c = h + mainArrayInterval
	}
	var bi = &BlockIndex{
		store: s,
		index: map[types.Hash]*types.Block{},
		main:  make([]*types.Block, c),
	}
	if highestBlock == nil {
		return bi, nil
	}

	blockHash := highestBlock.GetHash()
	for {
		block, err := (*s).GetBlock(blockHash.Bytes())
		if err != nil {
			return nil, err
		}
		bi.main[block.Height] = block
		bi.index[blockHash] = block

		if block.GetHash() == types.GenesisBlockHash() {
			log.Logger.Info("main chain and index complete.")
			break
		}
		//prev
		blockHash = block.PrevBlockHash
	}
	return bi, nil
}

func (bi *BlockIndex) exists(block *types.Block) bool {
	blockHash := block.GetHash()
	if _, ok := bi.index[blockHash]; ok {
		log.Logger.Infof("bock(hash=%x) exists in index map", blockHash)
		return true
	}
	if bi.main[block.Height] != nil {
		log.Logger.Infof("bock(hash=%x height=%d) exists in main chain", blockHash, block.Height)
		return true
	}
	return false
}

func (bi *BlockIndex) addMain(block *types.Block) (rollbackFn func()) {
	bi.maybeNeedExpansion()
	bi.main[block.Height] = block
	return func() {
		bi.main[block.Height] = nil
	}
}

func (bi *BlockIndex) addIndex(block *types.Block) (rollbackFn func()) {
	blockHash := block.GetHash()
	bi.index[blockHash] = block
	return func() {
		delete(bi.index, blockHash)
	}
}

func (bi *BlockIndex) maybeNeedExpansion() {
	len := len(bi.main)
	if len == cap(bi.main) {
		newArr := make([]*types.Block, len+mainArrayInterval)
		copy(newArr, bi.main)
		bi.main = newArr
	}
}

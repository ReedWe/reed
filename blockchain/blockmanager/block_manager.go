// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package blockmanager

import (
	"github.com/reed/blockchain/store"
	"github.com/reed/types"
	"sync"
)

type BlockManager struct {
	store        *store.Store
	HighestBlock *types.Hash
	BlockIndex   *BlockIndex
	mtx          sync.RWMutex
}

func NewBlockManager(s *store.Store) (*BlockManager, error) {
	highest, err := (*s).GetHighestBlock()
	if err != nil {
		return nil, err
	}

	index, err := NewBlockIndex(s, highest)
	if err != nil {
		return nil, err
	}

	return &BlockManager{
		store:        s,
		HighestBlock: highest,
		BlockIndex:   index,
	}, nil

}

func (bm *BlockManager) SaveNewBlock(block *types.Block) (exists bool, err error) {
	bm.mtx.Lock()
	defer bm.mtx.Unlock()

	if bm.BlockIndex.Exists(block) {
		return true, nil
	}
	bm.BlockIndex.AddBlock(block)

	if err := (*bm.store).AddBlock(block); err != nil {
		bm.BlockIndex.RollbackAddBlock(block)
		return false, err
	}
	hash := block.GetHash()
	bm.HighestBlock = &hash

	return false, nil
}

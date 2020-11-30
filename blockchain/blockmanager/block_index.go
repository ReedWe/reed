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
	index map[*types.Hash]*types.Block
	main  []*types.Block
}

func NewBlockIndex(s *store.Store, highestBlockHash *types.Hash) (*BlockIndex, error) {

	var bi = &BlockIndex{
		store: s,
	}

	blockHash := highestBlockHash
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

func (bi *BlockIndex) Exists(block *types.Block) bool {
	blockHash := block.GetHash()
	if bi.main[block.Height] != nil {
		log.Logger.Info("bock(hash=%x height=%d) exists in main chain", blockHash, block.Height)
		return true
	}
	if bi.index[&blockHash] != nil {
		log.Logger.Info("bock(hash=%x) exists in index", blockHash)
		return true
	}
	return false
}

func (bi *BlockIndex) AddBlock(block *types.Block) {
	blockHash := block.GetHash()
	bi.main[block.Height] = block
	bi.index[&blockHash] = block
}

func (bi *BlockIndex) RollbackAddBlock(block *types.Block) {
	blockHash := block.GetHash()
	bi.main[block.Height] = nil
	bi.index[&blockHash] = nil
}

// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package blockmanager

import (
	"fmt"
	"github.com/reed/blockchain/store"
	"github.com/reed/log"
	"github.com/reed/types"
	"sync"
)

//										/ [B1]				/ [B1]
//	[A]			[A]<--[B1]			[A]					[A]
//										\ [B2]				\ [B2]<--[C]
//
//	i[A]		i[A,B1]				i[A,B1,B2]			i[A,B1,B2,C]
//	m[A]		m[A,B1]				m[A,B1]				m[A,B2,C]
//
//	i:index
//	m:main

type BlockManager struct {
	store            *store.Store
	blockIndex       *BlockIndex
	highestBlock     *types.Block
	errs             map[types.Hash]error
	blockReceptionCh chan<- *types.RecvWrap
	mtx              sync.RWMutex
}

func NewBlockManager(s *store.Store, highestBlock *types.Block, rCh chan<- *types.RecvWrap) (*BlockManager, error) {
	index, err := NewBlockIndex(s, highestBlock)
	if err != nil {
		return nil, err
	}

	return &BlockManager{
		store:            s,
		blockIndex:       index,
		highestBlock:     highestBlock,
		blockReceptionCh: rCh,
	}, nil

}

func (bm *BlockManager) AddNewBlock(block *types.Block) (exists bool, err error) {
	bm.mtx.Lock()
	defer bm.mtx.Unlock()

	blockHash := block.GetHash()
	if _, ok := bm.errs[blockHash]; ok {
		log.Logger.Infof("block(hash=%x) exists in errs map", blockHash)
		return true, nil
	}
	if bm.blockIndex.exists(block) {
		return true, nil
	}

	if block.PrevBlockHash == types.DefHash() || block.PrevBlockHash == bm.highestBlock.GetHash() {
		amRollbackFn := bm.blockIndex.addMain(block)
		aiRollbackFn := bm.blockIndex.addIndex(block)
		if err := (*bm.store).SaveBlock(block); err != nil {
			bm.errs[blockHash] = err
			amRollbackFn()
			aiRollbackFn()
			return false, err
		}
		bm.highestBlock = block
	} else {
		// TODO reorganize not complete
		bm.blockIndex.addIndex(block)
		if block.Height > bm.highestBlock.Height && block.BigNumber.Cmp(&bm.highestBlock.BigNumber) == 1 {
			log.Logger.Infof("ready to reorganize")
			if err := bm.reorganize(block); err != nil {
				return false, err
			}
		}
	}
	return false, nil
}

func (bm *BlockManager) reorganize(block *types.Block) error {
	reserves, discards, err := bm.calcFork(block, nil)
	if err != nil {
		return err
	}

	//TODO remember set highestBlock
	//bm.highestBlock =
	fmt.Println(reserves, discards)

	//TODO SendBreakWork maybe false: block == newHighestBlock && config.Default.Mining
	//bm.blockReceptionCh <- &types.RecvWrap{Block: block, SendBreakWork: true}
	return nil
}

func (bm *BlockManager) calcFork(block *types.Block, highestBlock *types.Block) ([]*types.Block, []*types.Block, error) {
	var (
		reserves []*types.Block
		discards []*types.Block
	)
	subPoint := block
	mainPoint := highestBlock

	reserves = append(reserves, subPoint)
	mainHeight := mainPoint.Height

	for {
		i, ok := bm.blockIndex.index[subPoint.PrevBlockHash]
		if !ok {
			break
		}
		if i.Height != mainHeight {
			subPoint = i
			reserves = append(reserves, i)
		} else {
			m := bm.blockIndex.main[mainHeight]
			if i == m {
				break
			} else {
				subPoint = i
				mainPoint = m
				reserves = append(reserves, i)
				discards = append(discards, m)
				mainHeight--
			}
		}
	}

	if subPoint.PrevBlockHash != mainPoint.PrevBlockHash {
		log.Logger.Infof("sub chain longer but are orphans,waiting for parent.")
		return nil, nil, nil
	}
	return reserves, discards, nil
}

func (bm *BlockManager) GetAncestor(height uint64) *types.Block {
	bm.mtx.RLock()
	defer bm.mtx.RUnlock()
	return bm.blockIndex.main[height]
}

func (bm *BlockManager) HighestBlock() *types.Block {
	return bm.highestBlock
}

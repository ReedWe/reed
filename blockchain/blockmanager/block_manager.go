// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package blockmanager

import (
	"fmt"
	"github.com/reed/blockchain/store"
	"github.com/reed/errors"
	"github.com/reed/log"
	"github.com/reed/types"
	"sync"
)

//										/ [B1]				/ [B1]
//	[A]			[A]<--[B1]			[A]					[A]
//										\ [B2]				\ [B2]<--[C]
//
//	f[]			f[]					f[B2]				f[B1]
//	i[A]		i[A,B1]				i[A,B1,B2]			i[A,B1,B2,C]
//	m[A]		m[A,B1]				m[A,B1]				m[A,B2,C]
//
//	f:failed
//	i:index
//	m:main

type BlockManager struct {
	store      *store.Store
	blockIndex *BlockIndex
	failed     map[types.Hash]*types.Block
	errs       map[types.Hash]error
	mtx        sync.RWMutex
}

var (
	addNewBlockErr = errors.New("add new block error")
)

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
		store:      s,
		blockIndex: index,
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
	if _, ok := bm.failed[blockHash]; ok {
		log.Logger.Infof("block(hash=%x) exists in failed map", blockHash)
		return true, nil
	}
	if bm.blockIndex.exists(block) {
		return true, nil
	}

	highest, err := bm.HighestBlock()
	if err != nil {
		return false, err
	}
	if *block.PrevBlockHash == highest.GetHash() {
		amRollbackFn := bm.blockIndex.addMain(block)
		aiRollbackFn := bm.blockIndex.addIndex(block)
		if err := (*bm.store).SaveBlockAndUpdateHighest(block); err != nil {
			bm.errs[blockHash] = err
			amRollbackFn()
			aiRollbackFn()
			return false, err
		}
	} else {
		bm.blockIndex.addIndex(block)
		bm.failed[blockHash] = block
		if block.Height > highest.Height && block.BigNumber.Cmp(&highest.BigNumber) == 1 {
			log.Logger.Infof("ready to reorganize...")
			bm.reorganize(block)
		}
	}
	return false, nil
}

func (bm *BlockManager) reorganize(block *types.Block) error {
	reserves, discards, err := bm.calcFork(block)
	if err != nil {
		return err
	}

	fmt.Println(reserves, discards)
	return nil
}

func (bm *BlockManager) calcFork(block *types.Block) ([]*types.Block, []*types.Block, error) {
	var (
		reserves []*types.Block
		discards []*types.Block
	)
	subPoint := block
	mainPoint, err := bm.HighestBlock()
	if err != nil {
		return nil, nil, err
	}

	reserves = append(reserves, subPoint)
	mainHeight := mainPoint.Height

	for {
		i, ok := bm.blockIndex.index[*subPoint.PrevBlockHash]
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
				reserves = append(reserves, i)
				discards = append(discards, m)
				mainHeight--
			}
		}
	}

	if subPoint.PrevBlockHash != mainPoint.PrevBlockHash {
		//subs are Orphans
		return nil, nil, nil
	}
	return reserves, discards, nil
}

func (bm *BlockManager) HighestBlock() (*types.Block, error) {
	block, err := (*bm.store).GetHighestBlock()
	if err != nil {
		return nil, err
	}
	return block, nil
}

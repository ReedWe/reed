// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package validation

import (
	bm "github.com/reed/blockchain/blockmanager"
	"github.com/reed/consensus/pow"
	"github.com/reed/errors"
	"github.com/reed/types"
)

var (
	blockHeightErr     = errors.New("invalid block height")
	blockDiffErr       = errors.New("invalid block difficulty value")
	blockNonceErr      = errors.New("invalid block nonce value")
	blockParentHashErr = errors.New("invalid block prevBlockHash")
)

func ValidateBlockHeader(block *types.Block, prev *types.Block, bm *bm.BlockManager) error {
	if block.Height != prev.Height+1 {
		return errors.Wrapf(blockHeightErr, "prev height %d,cur height %d", prev.Height, block.Height)
	}

	difficulty := pow.GetDifficulty(block, bm.GetAncestor)
	if block.BigNumber.Cmp(&difficulty) != 0 {
		return errors.Wrap(blockDiffErr)
	}
	if !pow.CheckProofOfWork(difficulty, block.GetHash()) {
		return errors.Wrap(blockNonceErr)
	}
	if block.PrevBlockHash != prev.GetHash() {
		return errors.Wrapf(blockParentHashErr, "expect %x, actual %x", prev.GetHash(), &block.PrevBlockHash)
	}

	//TODO Version

	//TODO timestamp
	return nil
}

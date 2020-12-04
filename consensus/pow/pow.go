// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package pow

import (
	"encoding/hex"
	"github.com/reed/types"
	"math/big"
)

const (
	diffLimitHex = "000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
	//diffLimitHex   = "00000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
	targetTimespan = 14 * 24 * 60 * 60 // two weeks
	targetSpacing  = 10 * 60
)

func GetDifficulty(block *types.Block, getAncestor func(height uint64) *types.Block) big.Int {
	if block.Height == 1 || block.Height%DifficultyAdjustmentInterval() != 1 {
		return block.BigNumber
	}
	ancestor := getAncestor(block.Height - DifficultyAdjustmentInterval())
	return CalcDifficulty(ancestor.Timestamp, block)
}

func CheckProofOfWork(target big.Int, hash types.Hash) bool {
	var hashIntVal big.Int
	hashIntVal.SetBytes(hash.Bytes())
	return hashIntVal.Cmp(&target) == -1
}

//	Calculate a new difficulty
//	b:cur block
func CalcDifficulty(prevEpochBlockTime uint64, b *types.Block) big.Int {
	actualTimespan := b.Timestamp - prevEpochBlockTime

	if actualTimespan < targetTimespan/4 {
		actualTimespan = targetTimespan / 4
	}
	if actualTimespan > targetSpacing*4 {
		actualTimespan = targetSpacing * 4
	}

	var newDiff big.Int

	newDiff.Mul(&b.BigNumber, new(big.Int).SetUint64(actualTimespan))
	newDiff.Div(&newDiff, new(big.Int).SetUint64(targetTimespan))

	diffLimit := DifficultyLimit()
	if newDiff.Cmp(&diffLimit) == 1 {
		newDiff = diffLimit
	}

	return newDiff
}

func DifficultyLimit() big.Int {
	var n big.Int
	ds, _ := hex.DecodeString(diffLimitHex)
	n.SetBytes(ds)
	return n
}

func DifficultyAdjustmentInterval() uint64 {
	return targetTimespan / targetSpacing
}

// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package merkle

import (
	"github.com/reed/crypto"
	"github.com/reed/types"
)

func ComputeMerkleRoot(txs []*types.Tx) types.Hash {
	if len(txs) == 0 {
		return types.DefHash()
	}

	var trees [][]byte
	for _, tx := range txs {
		trees = append(trees, tx.GetID().Bytes())
	}

	if len(trees)&1 == 1 {
		trees = append(trees, trees[len(trees)-1])
	}

	count := len(trees)
	for {
		if count&1 == 1 {
			trees[count] = trees[count-1]
			count++
		}

		c := 0
		for i := 0; i < count; i++ {
			if i&1 == 1 {
				trees[i/2] = crypto.Sha256(trees[i-1], trees[i])
				c++
			}
		}
		if c == 1 {
			break
		}
		count = c
	}

	return types.BytesToHash(trees[0])
}

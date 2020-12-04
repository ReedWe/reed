// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package types

import "github.com/reed/blockchain/merkle"

type Block struct {
	BlockHeader
	Transactions []*Tx
}

func GetGenesisBlock() *Block {
	return &Block{
		BlockHeader:  *GetGenesisHeader(),
		Transactions: []*Tx{},
	}
}

func (b *Block) ComputeMerkleRootHash() {
	b.MerkleRootHash = merkle.ComputeMerkleRoot(b.Transactions)
}

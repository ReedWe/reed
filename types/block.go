// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package types

type Block struct {
	BlockHeader
	Transactions []*Tx
}

func NewBlock(prev *Block, Txs []*Tx) *Block {
	b := &Block{
		BlockHeader:  *prev.Copy(),
		Transactions: Txs,
	}
	return b
}

func GenesisBlock() *Block {
	return &Block{
		BlockHeader:  *GenesisHeader(),
		Transactions: []*Tx{},
	}
}

// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package types

import (
	"bytes"
	"github.com/reed/common/byteutil/byteconv"
	"github.com/reed/crypto"
	"math/big"
)

type BlockHeader struct {
	Height         uint64
	PrevBlockHash  *Hash
	MerkleRootHash *Hash
	Timestamp      uint64
	Nonce          uint64
	BigNumber      big.Int
}

func (bh *BlockHeader) GetHash() Hash {
	msg := bytes.Join([][]byte{
		byteconv.Uint64ToBytes(bh.Height),
		bh.PrevBlockHash.Bytes(),
		bh.MerkleRootHash.Bytes(),
		byteconv.Uint64ToBytes(bh.Timestamp),
		byteconv.Uint64ToBytes(bh.Nonce),
		bh.BigNumber.Bytes(),
	}, []byte{})

	return BytesToHash(crypto.Sha256(msg))
}

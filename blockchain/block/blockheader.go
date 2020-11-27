// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package block

import (
	"bytes"
	"github.com/reed/common/byteutil/byteconv"
	"github.com/reed/crypto"
	"github.com/reed/types"
	"math/big"
)

type Header struct {
	Height         uint64
	PrevBlockHash  *types.Hash
	MerkleRootHash *types.Hash
	Timestamp      uint64
	Nonce          uint64
	BigNumber      big.Int
}

func (bh *Header) GetHash() types.Hash {
	msg := bytes.Join([][]byte{
		byteconv.Uint64ToBytes(bh.Height),
		bh.PrevBlockHash.Bytes(),
		bh.MerkleRootHash.Bytes(),
		byteconv.Uint64ToBytes(bh.Timestamp),
		byteconv.Uint64ToBytes(bh.Nonce),
		bh.BigNumber.Bytes(),
	}, []byte{})

	return types.BytesToHash(crypto.Sha256(msg))
}

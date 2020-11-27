// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package block

import (
	"bytes"
	"encoding/binary"
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
	heightB := make([]byte, 8)
	tsB := make([]byte, 8)
	nonceB := make([]byte, 8)

	binary.LittleEndian.PutUint64(heightB, bh.Height)
	binary.LittleEndian.PutUint64(tsB, bh.Height)
	binary.LittleEndian.PutUint64(nonceB, bh.Height)

	msg := bytes.Join([][]byte{
		heightB,
		bh.PrevBlockHash.Bytes(),
		bh.MerkleRootHash.Bytes(),
		tsB,
		nonceB,
		bh.BigNumber.Bytes(),
	}, []byte{})

	return types.BytesToHash(crypto.Sha256(msg))
}

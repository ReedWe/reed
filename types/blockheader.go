// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package types

import (
	"bytes"
	"encoding/binary"
	"github.com/reed/crypto"
	"math/big"
)

type BlockHeader struct {
	Height      uint64
	PrevBlockId *Hash
	Timestamp   uint64
	Nonce       big.Int
	Bits        big.Int
}

func (bh *BlockHeader) GetHash() Hash {
	heightB := make([]byte, 8)
	tsB := make([]byte, 8)
	nonceB := make([]byte, 8)

	binary.LittleEndian.PutUint64(heightB, bh.Height)
	binary.LittleEndian.PutUint64(tsB, bh.Height)
	binary.LittleEndian.PutUint64(nonceB, bh.Height)

	msg := bytes.Join([][]byte{
		heightB,
		bh.PrevBlockId.Bytes(),
		tsB,
		nonceB,
		bh.Bits.Bytes(),
	}, []byte{})

	return BytesToHash(crypto.Sha256(msg))
}

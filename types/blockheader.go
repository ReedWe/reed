// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package types

import (
	"bytes"
	"github.com/reed/common/byteutil/byteconv"
	"github.com/reed/crypto"
	"math/big"
	"time"
)

type BlockHeader struct {
	Height         uint64
	PrevBlockHash  Hash
	MerkleRootHash Hash
	Timestamp      uint64
	Nonce          uint64
	BigNumber      big.Int
	Version        uint64
}

func (bh *BlockHeader) GetHash() Hash {
	msg := bytes.Join([][]byte{
		byteconv.Uint64ToByte(bh.Height),
		bh.PrevBlockHash.Bytes(),
		bh.MerkleRootHash.Bytes(),
		byteconv.Uint64ToByte(bh.Timestamp),
		byteconv.Uint64ToByte(bh.Nonce),
		bh.BigNumber.Bytes(),
		byteconv.Uint64ToByte(bh.Version),
	}, []byte{})

	return BytesToHash(crypto.Sha256(msg))
}

func (bh *BlockHeader) Copy() *BlockHeader {
	return &BlockHeader{
		Height:        bh.Height + 1,
		PrevBlockHash: bh.GetHash(),
		Timestamp:     uint64(time.Now().Unix()),
		Nonce:         0,
		BigNumber:     *big.NewInt(0),
		Version:       10000,
	}
}

func GetGenesisHeader() *BlockHeader {
	return &BlockHeader{
		Height:        1,
		PrevBlockHash: GenesisParentHash(),
		Timestamp:     uint64(time.Now().Unix()),
		Nonce:         0,
		BigNumber:     *new(big.Int),
		Version:       10000,
	}
}

// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package blockchain

const (
	initAmt = 16
)

func CalcCoinbaseAmt(curHeight uint64) uint64 {
	halving := curHeight / 210000 //TODO config params
	if halving >= 64 {
		return 0
	}
	amt := initAmt
	amt >>= halving
	return uint64(amt)
}

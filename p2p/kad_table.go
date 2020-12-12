// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package p2p

import (
	"fmt"
	"github.com/reed/types"
)

const (
	IDBits = len(types.Hash{}) * 8
)

type Table struct {
	Bucket [IDBits][]*Node
}

func NewTable() *Table {
	return &Table{
		Bucket: [IDBits][]*Node{},
	}
}

// logarithmDist return distance between a and b
// log2(a^b)

//	k-bucket	distance	description
//	0			[2^0,2^1)	存放距离为1，且前255bit相同，第256bit开始不同（即前255bit为0）
//	1			[2^1,2^2)	存放距离为2~3，且前254bit相同，第255bit开始不同
//	2			[2^2,2^3)	存放距离为4~7，且前253bit相同，第254bit开始不同
//	...
//	MEMO:
//	ID长度为32Byte，256bit。
//	上面循环每一位，进行异或（^）操作，结果0表示相同，1表示不同
//	所以“前导0个数为255”表示有255个bit是相同的
func LogarithmDist(a, b types.Hash) int {
	lz := 0
	for i := range a {
		x := a[i] ^ b[i]
		if x != 0 {
			fmt.Println("")
		}
		if x == 0 {
			lz += 8 // [0,0,0,0,0,0,0,0]
		} else {
			lz += lzcount[i]
		}
	}
	fmt.Println(lz)
	return len(a)*8 - lz
}

func Calc(a, b types.Hash) int {
	for i := range a {
		x := a[i] ^ b[i]
		if x != 0 {
			lz := i*8 + lzcount[x] //256bit leading zero counts
			return IDBits - 1 - lz
		}
	}
	return 0
}

// table of leading zero counts for bytes [0..255]
var lzcount = [256]int{
	8, 7, 6, 6, 5, 5, 5, 5,
	4, 4, 4, 4, 4, 4, 4, 4,
	3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3,
	2, 2, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 2, 2,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

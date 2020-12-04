// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package merkle

import (
	"bytes"
	"encoding/hex"
	"github.com/reed/types"
	"testing"
)

func TestComputeMerkleRoot(t *testing.T) {

	txs := []*types.Tx{
		{
			TxInput: []*types.TxInput{
				{
					ID: types.BytesToHash([]byte{1, 2}),
				},
			},
		},
		{
			TxInput: []*types.TxInput{
				{
					ID: types.BytesToHash([]byte{1, 2, 3, 4}),
				},
			},
		}, {
			TxInput: []*types.TxInput{
				{
					ID: types.BytesToHash([]byte{4, 5, 6}),
				},
			},
		}, {
			TxInput: []*types.TxInput{
				{
					ID: types.BytesToHash([]byte{7, 8, 9}),
				},
			},
		}, {
			TxInput: []*types.TxInput{
				{
					ID: types.BytesToHash([]byte{11, 12, 13}),
				},
			},
		},
		{
			TxInput: []*types.TxInput{
				{
					ID: types.BytesToHash([]byte{20, 21, 22}),
				},
			},
		},
	}

	l1H12, _ := hex.DecodeString("cc0be05e96991abedc8f594bf77aa3e17b8618239b2051181e04a13696619f06")
	//l1H34, _ := hex.DecodeString("ea6c1819a5424c12eb33694e36378cf9b5d962908d9b4d04fbd0d1ec7e963078")
	//l1H56, _ := hex.DecodeString("f73d146ecd2f264ca9666f5a053d421b6f5630a6be06c8c48c1463dc4c3ba9c2")
	//
	//l2H1234, _ := hex.DecodeString("f151753a4c533dd99b4b9cd562ee4b0e77645d22899c70dde351d8019e54ba65")
	//l2H5656, _ := hex.DecodeString("6358f22e18abf603e315c0976419cbc9aa96ec3c9078b8eac1463b9539c6096f")

	l3H123456, _ := hex.DecodeString("5f4216934e19cc5c5a114a277c9366c43dc52080163192dc299a296900652c80")

	if !bytes.Equal(l1H12, ComputeMerkleRoot(txs[:2]).Bytes()) {
		t.Error("compute two transactions merkle hash error.")
	}

	if !bytes.Equal(l3H123456, ComputeMerkleRoot(txs).Bytes()) {
		t.Error("compute transactions merkle hash error.")
	}

}

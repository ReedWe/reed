// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package types

type RecvBlockType int

const (
	RecvBlockTypeMined = iota
	RecvBlockTypeRemote
	RecvBlockTypeReorganize
)

type RecvWrap struct {
	SendBreakWork bool
	Block         *Block
}

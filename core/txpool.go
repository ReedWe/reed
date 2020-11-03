package core

import "github.com/tybc/core/types"

type Txpool struct {
	Txs map[types.Hash]*types.Tx
}

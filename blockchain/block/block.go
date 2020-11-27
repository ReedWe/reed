package block

import "github.com/reed/types"

type Block struct {
	Header
	Transactions *[]types.Tx
}



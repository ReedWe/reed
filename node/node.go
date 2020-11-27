// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package node

import (
	"github.com/reed/types"
	"github.com/tendermint/tmlibs/common"
)

type Node struct {
	common.BaseService
	api *types.API
}

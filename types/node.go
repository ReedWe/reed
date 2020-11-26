// Copyright 2020 The Tybc Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package types

import "github.com/tendermint/tmlibs/common"

type Node struct {
	common.BaseService
	api *API
}

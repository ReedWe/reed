// Copyright 2020 The Tybc Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package types

import "net/http"

type API struct {
	Chain  Chain
	Server *http.Server
}

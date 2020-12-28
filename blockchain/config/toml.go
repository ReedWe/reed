// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package config

import (
	"github.com/tendermint/tmlibs/common"
	"path"
)

var (
	v = `
# Reed config TOML document.
version=10000
homeDir=""
dataBasePath="database"
logPath="log"
logAge=86400
mining=false

apiAddr="0.0.0.0:9880"

ourNode="67032b2b262d837fbe2a0608409986c571350689@127.0.0.1:30398"
seeds="67032b2b262d837fbe2a0608409986c571350689@127.0.0.1:30398"
LockName="LOCK"
`
)

func GenerateConfigIfNotExist(dir string) {
	common.EnsureDir(dir, 0700)
	configFilePath := path.Join(dir, "config.toml")

	if !common.FileExists(configFilePath) {
		common.MustWriteFile(configFilePath, []byte(v), 0644)
	}
}

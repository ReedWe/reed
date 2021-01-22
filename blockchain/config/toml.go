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

ourNode="a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3@127.0.0.1:30398"
seeds="a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3@127.0.0.1:30398"
LockName="LOCK"
`
)

func GenerateConfigIfNotExist(dir, fileName string) string{
	common.EnsureDir(dir, 0700)
	configFilePath := path.Join(dir, fileName)

	if !common.FileExists(configFilePath) {
		common.MustWriteFile(configFilePath, []byte(v), 0644)
	}
	return configFilePath
}

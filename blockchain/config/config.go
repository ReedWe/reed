// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package config

import (
	"github.com/BurntSushi/toml"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

var (
	Default = &Config{}
)

type Config struct {
	basic
	node
	p2p
}

type basic struct {
	Version uint64
	HomeDir string

	DataBasePath string

	LogPath string
	LogAge  uint32

	Mining bool
}

type node struct {
	APIAddr string
}

type p2p struct {
	Seeds string

	OurNode  string
	LockName string
}

func init() {
	homeDir := ""
	home := os.Getenv("HOME")
	if home == "" {
		if u, err := user.Current(); err == nil {
			home = u.HomeDir
		}
	}
	switch runtime.GOOS {
	case "darwin":
		homeDir = filepath.Join(home, "Library", "Reed")
	case "windows":
		localappdata := os.Getenv("LOCALAPPDATA")
		if localappdata != "" {
			homeDir = filepath.Join(localappdata, "Reed")
		} else {
			homeDir = filepath.Join(home, "AppData", "Local", "Reed")
		}
	default:
		homeDir = filepath.Join(home, ".reed")
	}

	GenerateConfigIfNotExist(homeDir)

	if _, err := toml.DecodeFile("config.toml", &Default); err != nil {
		panic("Failed to decode config toml:" + err.Error())
	}
	Default.HomeDir = homeDir
}

func DatabaseDir() string {
	return rootify(Default.DataBasePath, Default.HomeDir)
}

func LogDir() string {
	return rootify(Default.LogPath, Default.HomeDir)
}

// helper function to make config creation independent of root dir
func rootify(path, root string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(root, path)
}

// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package config

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

var Default = Config{
	dataBasePath: "database",

	logPath: "log",
	LogAge:  60 * 60 * 24, //Sec
}

type Config struct {
	HomeDir      string
	dataBasePath string

	logPath string
	LogAge  uint32
}

func init() {
	home := os.Getenv("HOME")
	if home == "" {
		if user, err := user.Current(); err == nil {
			home = user.HomeDir
		}
	}
	switch runtime.GOOS {
	case "darwin":
		Default.HomeDir = filepath.Join(home, "Library", "Reed")
	case "windows":
		localappdata := os.Getenv("LOCALAPPDATA")
		if localappdata != "" {
			Default.HomeDir = filepath.Join(localappdata, "Reed")
		} else {
			Default.HomeDir = filepath.Join(home, "AppData", "Local", "Reed")
		}
	default:
		Default.HomeDir = filepath.Join(home, ".reed")
	}
}

func DatabaseDir() string {
	return rootify(Default.dataBasePath, Default.HomeDir)
}

func LogDir() string {
	return rootify(Default.logPath, Default.HomeDir)
}

// helper function to make config creation independent of root dir
func rootify(path, root string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(root, path)
}

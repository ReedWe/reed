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

var (
	Default = Config{
		Version:      10000,
		dataBasePath: "database2",

		logPath: "log",
		LogAge:  60 * 60 * 24, //Sec

		Mining: false,
		//
		//LocalAddr: "127.0.0.1:30398",
		//Seeds:     defaultSeed(),
		//OurID:     "67032b2b262d837fbe2a0608409986c571350689",
		//LockName:  "LOCK",
		//APIAddr:   ":9888",
		//
		LocalAddr: "127.0.0.1:30399",
		Seeds:     defaultSeed(),
		OurID:     "569188e0b7e1abdb9ac700fb97c3a4c3f749b2ea",
		LockName:  "LOCK2",
		APIAddr:   ":9889",
	}
)

type Config struct {
	Version      uint64
	HomeDir      string

	dataBasePath string

	logPath string
	LogAge  uint32

	Mining bool

	//P2P
	LocalAddr string
	Seeds     []string

	OurID    string
	LockName string
	APIAddr  string
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

func defaultSeed() []string {
	var ss []string
	ss = append(ss, "127.0.0.1:30399")
	return ss
}

// helper function to make config creation independent of root dir
func rootify(path, root string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(root, path)
}
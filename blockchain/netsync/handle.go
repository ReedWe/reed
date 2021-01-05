// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package netsync

const (
	getBlockHeaderCode = byte(0x01)
)

func Handle(msg []byte) ([]byte, error) {
	var writeBytes []byte
	switch msgType := msg[0]; msgType {
	case getBlockHeaderCode:
		writeBytes = []byte("this is my block header")
	}
	return writeBytes, nil
}

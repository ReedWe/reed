// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package netsync

import "fmt"

const (
	getBlockHeaderCode     = byte(0x01)
	getBlockHeaderRespCode = byte(0x02)
)

type HandleServ struct {
}

func NewHandleServ() *HandleServ {
	return &HandleServ{}
}

func (h HandleServ) Receive(msg []byte) []byte {
	var writeBytes []byte
	switch msgType := msg[0]; msgType {
	case getBlockHeaderCode:
		writeBytes = []byte("this is my block header")
	case getBlockHeaderRespCode:
		fmt.Println("receive a getBlockHeaderRespCode:")
		fmt.Printf(string(msg[1:]))
	}
	return writeBytes
}

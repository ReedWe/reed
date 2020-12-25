// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package crypto

import (
	"crypto/sha1"
)

func Sha1(data []byte) []byte {
	d := sha1.New()
	d.Write(data)
	return d.Sum(nil)
}

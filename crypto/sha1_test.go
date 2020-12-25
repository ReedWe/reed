// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package crypto

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestSha1(t *testing.T) {
	r, _ := hex.DecodeString("40bd001563085fc35165329ea1ff5c5ecbdbbeef")
	h := Sha1([]byte("123"))

	if !bytes.Equal(r, h) {
		t.Error("sha1 calculate error")
	}
}

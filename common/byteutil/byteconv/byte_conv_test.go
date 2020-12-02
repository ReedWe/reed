// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package byteconv

import (
	"fmt"
	"testing"
)

func TestUint64ToBytes(t *testing.T) {
	i := uint64(99999)
	bytes := Uint64ToByte(i)
	fmt.Println(bytes)

	i2 := ByteToUint64(bytes)
	fmt.Println(i2)

	if i != i2 {
		t.Error("convert error")
	}
}

package types

import (
	"bytes"
	"encoding/hex"
)

const (
	HashLength = 32
)

type Hash [HashLength]byte

func BytesToHash(b []byte) Hash {
	var h Hash
	h.SetBytes(b)
	return h
}

func DefHash() Hash {
	return BytesToHash([]byte{0})
}

// SetBytes sets the hash to the value of b.
// If b is larger than len(h), b will be cropped from the left.
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}

// Bytes gets the byte representation of the underlying hash.
func (h Hash) Bytes() []byte { return h[:] }

func (h Hash) ToString() string {
	return hex.EncodeToString(h.Bytes())
}

func (h Hash) HashEqual(b Hash) bool {
	return bytes.Equal(h.Bytes(), b.Bytes())
}

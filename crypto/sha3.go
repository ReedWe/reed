package crypto

import (
	"golang.org/x/crypto/sha3"
	"hash"
)

func Sha256(data ...[]byte) []byte {
	d := sha3.New256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}

func New256() hash.Hash {
	return sha3.New256()
}
func Sum(hash hash.Hash) []byte {
	return hash.Sum(nil)
}

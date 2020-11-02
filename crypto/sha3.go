package crypto

import (
	"golang.org/x/crypto/sha3"
)

func Sha256(data ...[]byte) []byte {
	d := sha3.New256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}

func Sha256T(data []byte) [32]byte {
	return sha3.Sum256(data)
}

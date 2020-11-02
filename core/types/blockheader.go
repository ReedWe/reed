package types

import "math/big"

type BlockHeader struct {
	Height      uint64
	PrevBlockId *[32]byte
	Timestamp   uint64
	Nonce       uint64
	Bits        big.Int
}

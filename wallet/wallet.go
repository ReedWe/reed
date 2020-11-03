package wallet

import (
	"encoding/hex"
	"golang.org/x/crypto/ed25519"
)

type Wallet struct {
	Pub ed25519.PublicKey
	Pk  ed25519.PrivateKey
}

func My() *Wallet {
	//hard code key
	pub, _ := hex.DecodeString("b12049d709358dc427433050625aa2135163181ccc320f22859d7c065ecc9dcb")
	pk, _ := hex.DecodeString("b4af5bb2bb1fb9086d2cac65a667f1810dfb0ddd904f2edc947227271fdcaba5b12049d709358dc427433050625aa2135163181ccc320f22859d7c065ecc9dcb")

	return &Wallet{
		Pub: pub,
		Pk:  pk,
	}
}

package crypto

import "crypto/ed25519"

func Sign(privateKey ed25519.PrivateKey, message []byte) []byte {
	return ed25519.Sign(privateKey, message)
}

package crypto

import "crypto/ed25519"

func Sign(privateKey ed25519.PrivateKey, message []byte) []byte {
	return ed25519.Sign(privateKey, message)
}

func Verify(publicKey ed25519.PublicKey, message, sig []byte) bool {
	return ed25519.Verify(publicKey, message, sig)
}

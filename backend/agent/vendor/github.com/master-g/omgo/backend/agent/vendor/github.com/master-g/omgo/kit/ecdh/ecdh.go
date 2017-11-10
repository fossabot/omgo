package ecdh

import (
	"crypto"
	"io"
)

// ECDH interface implement key exchange via Diffie-Hellman algorithm
type ECDH interface {
	GenerateECKey(io.Reader) (crypto.PrivateKey, crypto.PublicKey, error)
	Marshal(crypto.PublicKey) []byte
	Unmarshal([]byte) (crypto.PublicKey, bool)
	GenerateSharedSecret(crypto.PrivateKey, crypto.PublicKey) ([]byte, error)

	GenerateECKeyBuf(io.Reader) ([]byte, []byte)
	GenerateSharedSecretBuf([]byte, []byte) []byte
}

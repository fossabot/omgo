package ecdh

import (
	"crypto"
	"golang.org/x/crypto/curve25519"
	"io"
)

type curve25519ECDH struct {
}

// NewCurve25519ECDH creates a new curve25519ECDH instance
func NewCurve25519ECDH() ECDH {
	return &curve25519ECDH{}
}

// GenerateECKey will generate a private key and a public key for key exchange via Diffie-Hellman algorithm
// see detail at https://tools.ietf.org/id/draft-josefsson-tls-curve25519-02.html
func (e *curve25519ECDH) GenerateECKey(rand io.Reader) (crypto.PrivateKey, crypto.PublicKey, error) {
	var pub, priv [32]byte
	var err error

	_, err = io.ReadFull(rand, priv[:])
	if err != nil {
		return nil, nil, err
	}

	// the most significant bit (bit 254) is set
	priv[31] &= 0x7F
	priv[31] |= 0x40
	// and the three least significant bits are cleared
	priv[0] &= 0xF8

	curve25519.ScalarBaseMult(&pub, &priv)

	return &priv, &pub, nil
}

func (e *curve25519ECDH) Marshal(p crypto.PublicKey) []byte {
	pub := p.(*[32]byte)
	return pub[:]
}

func (e *curve25519ECDH) Unmarshal(data []byte) (crypto.PublicKey, bool) {
	var pub [32]byte
	if len(data) != 32 {
		return nil, false
	}

	copy(pub[:], data)
	return &pub, true
}

func (e *curve25519ECDH) GenerateSharedSecret(privateKey crypto.PrivateKey, publicKey crypto.PublicKey) ([]byte, error) {
	var priv, pub, secret *[32]byte

	priv = privateKey.(*[32]byte)
	pub = publicKey.(*[32]byte)
	secret = new([32]byte)

	curve25519.ScalarMult(secret, priv, pub)

	return secret[:], nil
}

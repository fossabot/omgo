package ecdh

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"encoding/hex"
	"log"
	"testing"
)

func TestCurve25519ECDH_GenerateSharedSecret(t *testing.T) {
	ecdh := NewCurve25519ECDH()
	testECDH(ecdh, t)
}

func testECDH(e ECDH, t testing.TB) {
	var privKey1, privKey2 crypto.PrivateKey
	var pubKey1, pubKey2 crypto.PublicKey
	var pubKey1Buf, pubKey2Buf []byte
	var err error
	var ok bool
	var secret1, secret2 []byte

	privKey1, pubKey1, err = e.GenerateECKey(rand.Reader)
	if err != nil {
		t.Error(err)
	}
	privKey2, pubKey2, err = e.GenerateECKey(rand.Reader)
	if err != nil {
		t.Error(err)
	}

	pubKey1Buf = e.Marshal(pubKey1)
	pubKey2Buf = e.Marshal(pubKey2)

	pubKey1, ok = e.Unmarshal(pubKey1Buf)
	if !ok {
		t.Fatalf("Unmarshal does not work")
	}

	pubKey2, ok = e.Unmarshal(pubKey2Buf)
	if !ok {
		t.Fatalf("Unmarshal does not work")
	}

	secret1, err = e.GenerateSharedSecret(privKey1, pubKey2)
	if err != nil {
		t.Error(err)
	}
	secret2, err = e.GenerateSharedSecret(privKey2, pubKey1)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(secret1, secret2) {
		t.Fatalf("The two shared keys: %d, %d do not match", secret1, secret2)
	}
}

func TestExchange(t *testing.T) {
	ecdh := NewCurve25519ECDH()
	cX1, cE1 := ecdh.GenerateECKeyBuf(rand.Reader)
	cX2, cE2 := ecdh.GenerateECKeyBuf(rand.Reader)

	sX1, sE1 := ecdh.GenerateECKeyBuf(rand.Reader)
	sX2, sE2 := ecdh.GenerateECKeyBuf(rand.Reader)

	sKey1 := ecdh.GenerateSharedSecretBuf(sX1, cE1)
	sKey2 := ecdh.GenerateSharedSecretBuf(sX2, cE2)

	cKey1 := ecdh.GenerateSharedSecretBuf(cX1, sE1)
	cKey2 := ecdh.GenerateSharedSecretBuf(cX2, sE2)

	if !bytes.Equal(sKey1, cKey1) || !bytes.Equal(sKey2, cKey2) {
		log.Println("sKey1:", hex.EncodeToString(sKey1))
		log.Println("sKey2:", hex.EncodeToString(sKey2))
		log.Println("cKey1:", hex.EncodeToString(cKey1))
		log.Println("cKey2:", hex.EncodeToString(cKey2))
		t.Fatalf("keys not match!")
	}
}

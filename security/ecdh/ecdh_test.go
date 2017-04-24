package ecdh

import (
	"bytes"
	"crypto"
	"crypto/rand"
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

	//fmt.Println(hex.EncodeToString(privKey1.(*[32]uint8)[:]))
	//fmt.Println(hex.EncodeToString(privKey2.(*[32]uint8)[:]))
	//fmt.Println(hex.EncodeToString(pubKey1.(*[32]uint8)[:]))
	//fmt.Println(hex.EncodeToString(pubKey2.(*[32]uint8)[:]))
	//fmt.Println(hex.EncodeToString(secret1))

	//p(privKey1.(*[32]uint8)[:])
	//p(privKey2.(*[32]uint8)[:])
	//p(pubKey1.(*[32]uint8)[:])
	//p(pubKey2.(*[32]uint8)[:])
	//p(secret1)
}

//func p(s []uint8) {
//	fmt.Printf("{")
//	for _, v := range s {
//		fmt.Printf("0x%02x, ", v)
//	}
//	fmt.Printf("}\n")
//}

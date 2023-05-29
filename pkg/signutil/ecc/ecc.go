package ecc

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"

	"github.com/huaxr/framework/pkg/toolutil"

	"github.com/ethereum/go-ethereum/crypto/ecies"
)

const split = "+"

var prv1 *ecdsa.PrivateKey
var prv2 *ecies.PrivateKey

func Encrypt(pt []byte) ([]byte, error) {
	ct, err := ecies.Encrypt(rand.Reader, &prv2.PublicKey, pt, nil, nil)
	return ct, err
}

func Decrypt(ct []byte) ([]byte, error) {
	pt, err := prv2.Decrypt(ct, nil, nil)
	return pt, err
}

func InitEcc(secret [40]byte) {
	var err error
	// Initialize the elliptic curve
	pubKeyCurve := elliptic.P256()
	// set secret can make you service distribute with one ecc Curve
	secBytes := bytes.NewBuffer(secret[:]) // io.rand
	prv1, err = ecdsa.GenerateKey(pubKeyCurve, secBytes)

	if err != nil {
		panic(err)
	}
	// Convert the standard package generation private key to the ECIES private key
	prv2 = ecies.ImportECDSA(prv1)
}

// sign feature
func Sign(pt []byte) (sign []byte, err error) {
	// Generate two big.ing according to the plaintext and private key
	r, s, err := ecdsa.Sign(rand.Reader, prv1, pt)
	if err != nil {
		return nil, err
	}
	rs, err := r.MarshalText()
	if err != nil {
		return nil, err
	}
	ss, err := s.MarshalText()
	if err != nil {
		return nil, err
	}
	// Merge R, s (split by "+") and return it as a signature
	var b bytes.Buffer
	b.Write(rs)
	b.Write([]byte(split))
	b.Write(ss)
	return b.Bytes(), nil
}

func SignVer(pt, sign []byte) bool {
	var (
		rInt, sInt big.Int
		err        error
	)
	// According to Sign, resolve R, s
	rs := bytes.Split(sign, []byte(split))
	err = rInt.UnmarshalText(rs[0])
	if err != nil {
		return false
	}
	err = sInt.UnmarshalText(rs[1])
	if err != nil {
		return false
	}
	// Verify the signature against the public key, plaintext, R, S
	v := ecdsa.Verify(&prv1.PublicKey, pt, &rInt, &sInt)
	return v
}

func Str40ToBytes40(st string) [40]byte {
	if len(st) != 40 {
		st = toolutil.GetRandomString(40)
	}
	var res [40]byte
	for i, j := range st {
		res[i] = byte(j)
	}
	return res
}

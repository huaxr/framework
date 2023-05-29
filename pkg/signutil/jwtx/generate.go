// Author: huaxinrui@tal.com
// Time:   2021/9/8 下午12:58
// Git:    huaxr

package jwtx

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/huaxr/framework/logx"

	"os"
)

const (
	PRIVATE = "/tmp/private.pem"
	PUBLIC  = "/tmp/public.pem"
	bits    = 1024
)

func Gen() {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		logx.L().Errorf("privateKey err %v", err)
		return
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	file, err := os.Create(PRIVATE)
	if err != nil {
		logx.L().Errorf("privateKey err %v", err)
		return
	}
	err = pem.Encode(file, block)
	if err != nil {
		logx.L().Errorf("privateKey err %v", err)
		return
	}

	publicKey := &privateKey.PublicKey
	defPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		logx.L().Errorf("privateKey err %v", err)
		return
	}
	block = &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: defPkix,
	}
	file, err = os.Create(PUBLIC)
	if err != nil {
		logx.L().Errorf("privateKey", "%v", err)
		return
	}
	err = pem.Encode(file, block)
	if err != nil {
		logx.L().Errorf("privateKey", "%v", err)
		return
	}
}

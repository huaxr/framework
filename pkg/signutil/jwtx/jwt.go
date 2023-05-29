package jwtx

import (
	"fmt"
	"time"

	"github.com/huaxr/framework/pkg/toolutil"
	"github.com/dgrijalva/jwt-go"
)

// todo: key escrow
const private = `
-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQC2LW4Th7ArAizmbncBlVhVlry7HSmfwcxyMltnXOAMOI3yR9Sq
PIw0Lrb1QikDWmJ3XHCQ1Dcq3wJxxhQ6gCHgyfs2F+plX7H/uJPSlsGZ5I/yRlqa
Gzf2p9vg/r6TJ3o/5J4SaRQvOwfYhbGFuLDvJGkeRGt8nU3K3wAvhJufZQIDAQAB
AoGBAIUv138Nv1ziHUNmVTjiH4+LQXWmz2yNudNvP2Xk/6PPoO8VVsQSugnYcUgD
U4qxBLXw7hbkH2UHX3kgcF+Il7rLt3ZLQSVTssBVUSFbKsDkC++nuLsINSEbUH2l
OxKDpOfrKVcfrLTr7AQ1Og00eHFtjgyA9jbk3KRkB6ffY/6tAkEAxQ0Lq2hNIeSG
3lC8B9FZLqEIamb+JL7DcSVSZSyhxgPusi/vD4vyNjlSNQn8ZZCv8Ba/lJsj82UZ
FIun7QS6/wJBAOytT+a7kl95UcNAp0h4Bh1oNPp+3ZjESKcjQBB5oNCa+ltIFAf+
SKcK3jY0gneCKeAEX0egY2BExPU/Rr0EmZsCQFjLVCLdUUSgkhXEE7cCI0nbzssD
tiogvDlUNBjbT9rHEtzAtN0wlujQU7cK1O1/kYiC97mjX0PinrafaABqTUkCQQCR
SNpxgtcZcHm2Z+vIWpU2XA+ZbWNOMb9/ie37rw3+wAPLIPXa6kdi8xLxJ06nWemm
sEhkyZn3MH/PJGaiBT7dAkEAgrWmIQmRrk+qaRxgUMDWFNUJgiPRzKuvFDr8YRYB
6lGjvJROZomBd8iF5/laMi8ozBEfS0h7za5zcLx+ct+MgA==
-----END RSA PRIVATE KEY-----`

const public = `
-----BEGIN RSA PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC2LW4Th7ArAizmbncBlVhVlry7
HSmfwcxyMltnXOAMOI3yR9SqPIw0Lrb1QikDWmJ3XHCQ1Dcq3wJxxhQ6gCHgyfs2
F+plX7H/uJPSlsGZ5I/yRlqaGzf2p9vg/r6TJ3o/5J4SaRQvOwfYhbGFuLDvJGke
RGt8nU3K3wAvhJufZQIDAQAB
-----END RSA PUBLIC KEY-----`

var (
	privateKey = toolutil.String2Byte(private)
	publicKey  = toolutil.String2Byte(public)
)

func GenTokenString(info map[string]interface{}) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return "", fmt.Errorf("parse private key err: %v\n", err)
	}
	token := jwt.New(jwt.SigningMethodRS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(24)).Unix()
	claims["iat"] = time.Now().Unix()
	for k, v := range info {
		claims[k] = v
	}
	token.Claims = claims

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

var keyFunc = func(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("unknown sign method: %v", token.Header["alg"])
	}
	key, _ := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	return key, nil
}

func CheckTokenString(token string) (*jwt.Token, error) {
	t, err := jwt.Parse(token, keyFunc)
	return t, err
}

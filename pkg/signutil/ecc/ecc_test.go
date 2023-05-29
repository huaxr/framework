package ecc

import (
	"testing"

	"github.com/huaxr/framework/pkg/toolutil"
)

// when use curve crypto, please define your curve coordinate first!
// using length of 40 strings to identify your X&Y
var x = Str40ToBytes40("2z4qr29x7ku1expf28is244efl2a0ok3ztpjesbv")
var y = Str40ToBytes40("2z4qr29x7ku1expf28is244efl2a0ok3ztpjesbw")

func TestEcc(t *testing.T) {
	InitEcc(x)
	b, err := Encrypt(toolutil.String2Byte("abcd"))
	if err != nil {
		t.Log(err)
		return
	}

	InitEcc(y)
	res, err := Decrypt(b)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(string(res))
}

func TestEccSign(t *testing.T) {
	InitEcc(x)
	pt := toolutil.String2Byte("Hello world @#!~%^&*()")
	sig, err := Sign(pt)
	if err != nil {
		t.Log(err)
		return
	}

	ok := SignVer(pt, sig)
	t.Log(ok)

	ok = SignVer(pt, append(sig, byte(1)))
	t.Log(ok)

	InitEcc(y)
	ok = SignVer(pt, sig)
	t.Log(ok)

	InitEcc(x)
	ok = SignVer(pt, sig)
	t.Log(ok)
}

func BenchmarkEccSign(b *testing.B) {
	InitEcc(x)
	pt := toolutil.String2Byte("abcdEFGHIGKLMNOPQRSTUvwxyz")

	b.Run("ECC", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sig, _ := Sign(pt)
			SignVer(pt, sig)
		}
	})

	b.StopTimer()
}

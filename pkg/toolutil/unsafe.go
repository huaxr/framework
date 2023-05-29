// Author: huaxinrui@tal.com
// Time:   2021/7/15 下午2:32
// Git:    huaxr

package toolutil

import (
	"reflect"
	"unsafe"
)

// zero copy using the underlying structure of slice&string
// unsafe pointer is a pointer reference associate with uintptr.
func String2Byte(s string) []byte {
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: stringHeader.Data,
		Len:  stringHeader.Len,
		Cap:  stringHeader.Len,
	}
	//runtime.KeepAlive(&s)
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func Bytes2string(b []byte) string {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))

	sh := reflect.StringHeader{
		Data: sliceHeader.Data,
		Len:  sliceHeader.Len,
	}

	//runtime.KeepAlive(&b)
	return *(*string)(unsafe.Pointer(&sh))
}

// GetMapSize return the substructure h-map first field length.
func GetMapSize(m map[string]interface{}) int {
	return **(**int)(unsafe.Pointer(&m))
}

func Float64bits(f float64) uint64 {
	return *(*uint64)(unsafe.Pointer(&f))
}

func Float64FromBits(b uint64) float64 {
	return *(*float64)(unsafe.Pointer(&b))
}

// Author: huaxr
// Time:   2021/7/15 下午2:34
// Git:    huaxr

package toolutil

import (
	"fmt"
	"log"
	"testing"
	"unsafe"
)

func TestGetMapSize(t *testing.T) {
	a := map[string]interface{}{"a": "b"}
	x := GetMapSize(a)
	t.Logf("%v", x)

	b := make(map[string]interface{})
	b["a"] = "c"
	b["x"] = "a"
	y := GetMapSize(b)
	t.Logf("%v", y)
}

func TestUnsafe(t *testing.T) {
	var x struct {
		a int64
		b bool
		c string
	}
	const M, N = unsafe.Sizeof(x.c), unsafe.Sizeof(x)
	fmt.Println(M, N) // 16 32

	fmt.Println(unsafe.Alignof(x.a)) // 8
	fmt.Println(unsafe.Alignof(x.b)) // 1
	fmt.Println(unsafe.Alignof(x.c)) // 8

	fmt.Println(unsafe.Offsetof(x.a)) // 0
	fmt.Println(unsafe.Offsetof(x.b)) // 8
	fmt.Println(unsafe.Offsetof(x.c)) // 16
}

func TestSafe2(t *testing.T) {
	type T struct {
		c string
	}
	type S struct {
		b bool
	}
	var x struct {
		a int64
		*S
		T
	}

	fmt.Println(unsafe.Offsetof(x.a)) // 0

	fmt.Println(unsafe.Offsetof(x.S)) // 8
	fmt.Println(unsafe.Offsetof(x.T)) // 16

	// 此行可以编译过，因为选择器x.c中的隐含字段T为非指针。
	fmt.Println(unsafe.Offsetof(x.c)) // 16

	// 此行编译不过，因为选择器x.b中的隐含字段S为指针。
	//fmt.Println(unsafe.Offsetof(x.b)) // error

	// 此行可以编译过，但是它将打印出字段b在x.S中的偏移量.
	fmt.Println(unsafe.Offsetof(x.S.b)) // 0
}

func TestOffset(ts *testing.T) {
	type T struct {
		x bool
		y [3]int16
	}

	const N = unsafe.Offsetof(T{}.y)
	const M = unsafe.Sizeof(T{}.y[0])

	t := T{y: [3]int16{123, 456, 789}}
	p := unsafe.Pointer(&t)
	// "uintptr(p) + N + M + M"为t.y[2]的内存地址。
	ty2 := (*int16)(unsafe.Pointer(uintptr(p) + N + M + M))
	fmt.Println(*ty2) // 789

	// 不应该如下拆开写，虽然这里都一样
	// ty2 := (*int16)(unsafe.Pointer(uintptr(p)+N+M+M))
	addr := uintptr(p) + N + M + M
	// ...（一些其它操作）
	// 从这里到下一行代码执行之前，t值将不再被任何值
	// 引用，所以垃圾回收器认为它可以被回收了。一旦
	// 它真地被回收了，下面继续使用t.y[2]值的曾经
	// 的地址是非法和危险的！另一个危险的原因是
	// t的地址在执行下一行之前可能改变（见事实三）。
	// 另一个潜在的危险是：如果在此期间发生了一些
	// 操作导致协程堆栈大小改变的情况，则记录在addr
	// 中的地址将失效。
	ty2 = (*int16)(unsafe.Pointer(addr))
	fmt.Println(*ty2)
}

// *unsafe.Pointer是一个类型安全指针类型
func SafePointer() {
	x := 123                // 类型为int
	p := unsafe.Pointer(&x) // 类型为unsafe.Pointer
	pp := &p                // 类型为*unsafe.Pointer
	p = unsafe.Pointer(pp)
	pp = (*unsafe.Pointer)(p)
}

func TestStr(t *testing.T) {
	type x struct {
		Name string
		Age  string
		Info map[string]interface{}
		A    *x
		O    int8
		I    []byte
	}
	sl := x{"hua", "rui", map[string]interface{}{"a": 1}, nil, 55, []byte{1}}

	const A = unsafe.Offsetof(x{}.Name)
	const N = unsafe.Offsetof(x{}.Age)
	const X = unsafe.Offsetof(x{}.Info)
	const XX = unsafe.Offsetof(x{}.A)
	const M = unsafe.Offsetof(x{}.O)
	const I = unsafe.Offsetof(x{}.I)

	log.Println(A, N, X, XX, M, I)
	P := uintptr(unsafe.Pointer(&sl))
	b := *(*int)(unsafe.Pointer(P + M))

	t.Logf("%v ", b)
}

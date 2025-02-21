package xstruct

import (
	"fmt"
	"unsafe"
)

// bad
type S1 struct {
	A int8
	B int64
	C int16
}

// [ A ] [ P P P P P P P ] [ B B B B B B B B ] [ C C ] [ P P ]
// A 占1字节
// P 填充7字节，为了保证B的初始地址是8的倍数
// B 8 字节
// C 2 字节
// P 为了保证整体结构是8的倍数，填充2字节

// good
type S2 struct {
	A int8  // 1
	C int16 // 2
	B int64 // 8
}

// [ A ] [P] [C C] [P P P P] [B B B B B B B B]
// A 1 字节
// P 7 字节，原因同上
// C 2 字节
// P 6 字节，保证B初始地址是8倍数
// B 8 字节
// 1 + 1 + 2 + 4 + 8 // 16

type S3 struct {
	B int64
	C int16
	A int8
}

// [B B B B B B B B] [C C] [P P P P P] [A]

func PrintSize() {
	fmt.Println(unsafe.Sizeof(S1{}))
	fmt.Println(unsafe.Sizeof(S2{}))
	fmt.Println(unsafe.Sizeof(S3{}))
}

type Inner struct {
	a int8  // 1
	b int32 // 4
	c int64 // 8
}

// [a] [P P P] [b b b b] [c c c c c c c c] 16

type Outer struct {
	x int16
	y Inner
	z int64
}

// [x x] [a] [p] [b b b b] [z z z z z z z z] 16
// 暂时不纠结该问题。
type X struct {
	A int    // 8
	B [0]int //
}

type Y struct {
	B [0]int // 0
	A int    // 8
}

func main() {
	var x X
	var y Y
	fmt.Println(unsafe.Sizeof(x))
	fmt.Println(unsafe.Sizeof(y))
}

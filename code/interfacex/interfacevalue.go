package interfacex

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// 接口赋值
func TestInterface() {
	var w io.Writer
	w = os.Stdout
	w = new(bytes.Buffer)
	w = nil

	fmt.Println(w)
}

func TestInterface2() {
	var a interface{} = 10
	var b interface{} = 10
	var c interface{} = "hello"
	var s1 interface{} = []int{1, 2}
	var s2 interface{} = []int{1, 2}

	fmt.Println(a == b)   // true (动态类型都是 int，值也相等)
	fmt.Println(a == c)   // false (动态类型不同)
	fmt.Println(s1 == s2) // 运行时 panic：runtime error: comparing uncomparable type []int

	////结构体包含接口的例子
	//type MyStruct struct{
	//	Val interface{}
	//}
	//s3 := MyStruct{Val: 10}
	//s4 := MyStruct{Val: 10}
	//s5 := MyStruct{Val: "hello"}
	//s6 := MyStruct{Val: []int{1,2}}
	//s7 := MyStruct{Val: []int{1,2}}
	//fmt.Println(s3 == s4)  // true
	//fmt.Println(s3 == s5) // false
	////fmt.Println(s6 == s7)  // 运行时 panic
	//
	//var x *int
	//var y *int
	//var ix interface{} = x
	//var iy interface{} = y
	//fmt.Println(ix == iy) //true, 都是nil指针，动态类型都是*int
	//
	//var z int
	//var iz interface{} = z
	//fmt.Println(ix == iz) // panic: 动态类型不同
}

// 注意，含有空指针的非空接口的情况
// 即接口的type不为空，值语义部分为空；
func TestNilInterface() {
	var p *int
	var i interface{} = p

	fmt.Println(i == nil)

	if i != nil {
		fmt.Println("i is not nil, even though p is nil")
	}
}

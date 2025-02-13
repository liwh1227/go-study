package mn

import "fmt"

func ExampleMakeNew() {
	var s = new([]string)
	var ss = make([]string, 0)

	fmt.Println("s:", s)                 // 输出: s: &[]
	fmt.Println("ss:", ss)               // 输出: ss: []
	fmt.Println("s == nil:", s == nil)   // 输出: s == nil: false
	fmt.Println("*s == nil:", *s == nil) // 输出: *s == nil: true
	fmt.Println("ss == nil:", ss == nil) // 输出: ss == nil: false

	// 使用 s (需要解引用)
	*s = append(*s, "hello")
	fmt.Println("s:", s)   // 输出: s: &[hello]
	fmt.Println("*s:", *s) // 输出: *s: [hello]

	// 使用 ss (直接使用)
	ss = append(ss, "world")
	fmt.Println("ss:", ss) // 输出: ss: [world]
	fmt.Println("len(ss):", len(ss))
	fmt.Println("cap(ss):", cap(ss))
}

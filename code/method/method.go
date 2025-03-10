package method

import (
	"fmt"
	"net/url"
)

type Point struct {
	X, Y int
}

func (p Point) Set(x, y int) {
	p.X = x
	p.Y = y
}

func (p *Point) Set2(x, y int) {
	p.X = x
	p.Y = y
}

func TestMethod2() {
	p := Point{
		X: 0,
		Y: 0,
	}

	p.Set2(100, 100)

	fmt.Println(p.X, p.Y)
}

func TestMethod() {
	m := url.Values{"lang": {"en"}}

	m.Add("item", "1")
	m.Add("item", "2")

	fmt.Printf("%v", m)
}

type Ss []int

func (s Ss) Append(x int) {
	fmt.Printf("Before append - Inside Append: len(s) = %d, cap(s) = %d, s = %v\n", len(s), cap(s), s)
	s = append(s, x) // 修改的是副本 s
	fmt.Printf("After append - Inside Append: len(s) = %d, cap(s) = %d, s = %v\n", len(s), cap(s), s)
}

type Values map[string][]string

func TestMethod3() {
	// 情况 1: 容量足够，不扩容
	ss1 := make(Ss, 0, 5)   // 创建一个长度为 0，容量为 5 的切片
	ss1 = append(ss1, 1, 2) // 先添加两个元素, 使其len=2
	fmt.Printf("Before Append - Outside: len(ss1) = %d, cap(ss1) = %d, ss1 = %v\n", len(ss1), cap(ss1), ss1)
	ss1[0] = 10000
	ss1.Append(3)
	fmt.Printf("After Append - Outside: len(ss1) = %d, cap(ss1) = %d, ss1 = %v\n", len(ss1), cap(ss1), ss1)

	fmt.Println("Underlying array:", ss1[:cap(ss1)]) // 打印底层数组

	fmt.Println("-----")

	// 情况 2: 容量不足，扩容
	ss2 := make(Ss, 0, 2)
	ss2 = append(ss2, 1, 2)
	fmt.Printf("Before Append - Outside: len(ss2) = %d, cap(ss2) = %d, ss2 = %v\n", len(ss2), cap(ss2), ss2)
	ss2.Append(3)
	fmt.Printf("After Append - Outside: len(ss2) = %d, cap(ss2) = %d, ss2 = %v\n", len(ss2), cap(ss2), ss2)
	fmt.Println("Underlying array:", ss2[:cap(ss2)])
}

func (v Values) Add(key, value string) {
	v[key] = append(v[key], value)
}

func TestMethod4() {
	vv := Values{
		"item": {"1900"},
	}
	vv.Add("item", "100")

	fmt.Println(vv)
}

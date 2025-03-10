package interfacex

import "fmt"

// 定义一个接口
type Shape interface {
	Area() float64
}

// 定义两个结构体，分别表示圆形和矩形
type Circle struct {
	Radius float64
}

type Rectangle struct {
	Width  float64
	Height float64
}

// Circle 类型实现了 Shape 接口的 Area 方法
func (c Circle) Area() float64 {
	return 3.14159 * c.Radius * c.Radius
}

// Rectangle 类型实现了 Shape 接口的 Area 方法
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// 这个函数接受一个 Shape 接口类型的参数
func printArea(s Shape) {
	fmt.Println("Area:", s.Area())
}

func TestShape() {
	c := Circle{Radius: 5}
	r := Rectangle{Width: 4, Height: 6}

	// 可以将 Circle 和 Rectangle 类型的变量传递给 printArea 函数
	// 因为它们都实现了 Shape 接口
	printArea(c) // 输出: Area: 78.53975
	printArea(r) // 输出: Area: 24
}

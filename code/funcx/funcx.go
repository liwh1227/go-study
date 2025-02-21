package funcx

import (
	"fmt"
)

func square() func() int {
	var x int
	// 返回匿名函数
	return func() int {
		x++
		return x * x
	}
}

func TestF() {
	f := square()
	fmt.Println(f())
	fmt.Println(f())
	fmt.Println(f())
	fmt.Println(f())
}

// 创建一个带有内部状态的函数
func createGreeter(greeting string) func(name string) string {
	return func(name string) string {
		return greeting + ", " + name + "!"
	}
}

func TestA() {
	sayHello := createGreeter("Hello")
	sayGoodBye := createGreeter("Goodbye")

	fmt.Println(sayHello("Alice"))
	fmt.Println(sayHello("Bob"))
	fmt.Println(sayGoodBye("Charlie"))
}

func printValues(values ...int) {
	for _, val := range values {
		func() {
			fmt.Println(val) // 闭包捕获了循环变量 val
		}() // 这里是定义并立即执行
	}
}

func TestC() {
	printValues(1, 2, 3)
}

// error example:
//func TestB() {
//	var rmdis []func()
//	for _, dir := range tempDirs() {
//		os.Mkdir(dir, 0755)
//		rmdis = append(rmdis, func() {
//			os.Remove(dir)
//		})
//	}
//}

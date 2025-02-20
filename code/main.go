package main

import "fmt"

const (
	a = iota
	b
	c = 10
	d
	_ = iota
	f
)

func main() {
	fmt.Println(a, b, c, d, f)
}

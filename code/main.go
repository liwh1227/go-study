package main

import (
	"fmt"
	"go-study/array"
)

func main() {
	var s = []int{1, 2, 3}

	fmt.Println(len(s), cap(s))

	s = array.AppendInt(s, 100)

	fmt.Println(len(s), cap(s))

	fmt.Printf("%v", s)
}

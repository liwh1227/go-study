package main

import "fmt"

func main() {
	var s = []int{
		0, 1, 2, 3, 4, 5,
	}
	//singleRoutine(s)

	//multiRoutine(s)
	//
	multiRoutine2(s)

	select {}
}

// 单协程打印
func singleRoutine(s []int) {

	for i := range s {
		fmt.Printf("%d ", s[i])
	}
}

func multiRoutine(s []int) {
	for i := range s {
		go func() {
			fmt.Printf("%d ", s[i])
		}()
	}
}

func multiRoutine2(s []int) {
	for i := range s {
		go func(val int) {
			fmt.Printf("%d ", val)
		}(s[i])
	}
}

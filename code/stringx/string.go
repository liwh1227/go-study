package stringx

import "fmt"

func Print() {
	rangeString()
}

func printString() {
	fmt.Println(`hello \n world`)
	fmt.Println("hello \nworld")
}

func rangeString() {
	var s = "hello,世界"
	// 输出：12
	// 其反映的是底层[]byte数组的长度，
	fmt.Println(len(s))

	for i, char := range s {
		// 如果按照下标打印，可能会出现乱码
		fmt.Printf("%c, %c\n", s[i], char)
	}

	// 使用rune打印
	// 输出：8
	// 其反映的是 unicode 码点的长度，[]rune(s)，通过utf8编码后，实际是utf8编码后的数组；
	runes := []rune(s)

	fmt.Println(len(runes))
	for i := 0; i < len(runes); i++ {
		fmt.Printf("index: %d, character %c\n", i, runes[i])
	}
}

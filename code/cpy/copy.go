package cpy

import "fmt"

func TestCopy() {
	// 示例 1：基本复制
	src1 := []int{1, 2, 3, 4, 5}
	dst1 := make([]int, 3) // 长度为 3，容量默认为 3

	copied1 := copy(dst1, src1)
	fmt.Println("Copied:", copied1) // 输出: Copied: 3
	fmt.Println("dst1:", dst1)      // 输出: dst1: [1 2 3]
	fmt.Println("src1:", src1)      // 输出: src1: [1 2 3 4 5] (src1 不变)

	// 示例 2：目标切片比源切片长
	src2 := []int{1, 2}
	dst2 := make([]int, 5)

	copied2 := copy(dst2, src2)
	fmt.Println("Copied:", copied2) // 输出: Copied: 2
	fmt.Println("dst2:", dst2)      // 输出: dst2: [1 2 0 0 0] (只复制了前两个元素)

	// 示例 3：源切片比目标切片长
	src3 := []int{1, 2, 3, 4, 5}
	dst3 := []int{0}

	copied3 := copy(dst3, src3)
	fmt.Println("Copied:", copied3) // 输出: Copied: 2
	fmt.Println("dst3:", dst3)      // 输出: dst3: [1 2] (只复制了前两个元素)

	// 示例 4：重叠切片
	arr := []int{1, 2, 3, 4, 5}
	src4 := arr[0:3] // [1 2 3]
	dst4 := arr[2:5] // [3 4 5]
	// dst4 和 src4 共享底层数组 arr，并且有重叠部分

	copied4 := copy(dst4, src4)
	fmt.Println("Copied:", copied4) // 输出: Copied: 3
	fmt.Println("dst4", dst4)
	fmt.Println("arr:", arr) // 输出: arr: [1 2 1 2 3] (注意 arr 的变化)

	// 示例 5: 将字符串复制到 []byte 切片
	str := "hello"
	bytes := make([]byte, len(str))

	copied5 := copy(bytes, str)
	fmt.Println("Copied:", copied5)            // 输出：Copied: 5
	fmt.Println("bytes:", bytes)               // 输出：bytes: [104 101 108 108 111] (ASCII 码)
	fmt.Printf("bytes as string: %s\n", bytes) // 输出：bytes as string: hello
}

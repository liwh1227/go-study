package array

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func Rune() {
	var runes []rune
	for _, r := range "Hello, 世界" {
		runes = append(runes, r)
	}
	fmt.Printf("%q\n", runes)
}

// slice 定义的 appendInt 方法
func AppendInt(x []int, y int) []int {
	// 定义一个新的slice
	var z []int
	zlen := len(x) + 1
	if zlen <= cap(x) {
		// slice 仍有增长空间，未达到 x 底层数组的cap
		// 将x的内容复制到z上
		z = x[:zlen]
	} else {
		// slice 无空间
		zcap := zlen
		if zcap < 2*len(x) {
			zcap = 2 * len(x)
		}
		z = make([]int, zlen, zcap)
		copy(z, x)
	}
	// 将last index赋值为y
	z[len(x)] = y
	return z
}

// 去除slice中的空串
func Nonempty(strings []string) []string {
	i := 0
	for _, s := range strings {
		if s != "" {
			strings[i] = s
			i++
		}
	}
	return strings[:i]
}

func Reverse(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func ReverseByPtr(s *[]int) {
	n := len(*s)

	left, right := 0, n-1
	for left < right {
		(*s)[left], (*s)[right] = (*s)[right], (*s)[left]
		left++
		right--
	}
}

/*
1. reverse (反转):
目的: 将序列中的元素 颠倒顺序。
效果: 第一个元素变成最后一个，第二个元素变成倒数第二个，依此类推。
示例:
原始序列：[1, 2, 3, 4, 5]
reverse 后：[5, 4, 3, 2, 1]

2. rotate (旋转):
目的: 将序列中的元素 循环移动 指定的位置。
效果:
向左旋转 (rotate left): 将序列的 前 k 个元素 移动到序列的 末尾。
向右旋转 (rotate right): 将序列的 后 k 个元素 移动到序列的 开头。
示例:
原始序列：[1, 2, 3, 4, 5]
向左旋转 2 位：[3, 4, 5, 1, 2]
向右旋转 2 位：[4, 5, 1, 2, 3]
*/
// 注意，k可能是负数，左转的负数意味着右转
// 例如，k = -2 ,意味着右转2位，相应的对于左转而言，就是左转3位
// -2 % 5 = -2，k = 3
func RotateLeftPlace(s []int, k int) {
	n := len(s)
	if n == 0 {
		return
	}
	// 保证k的范围是在0～n-1
	// 从数学上讲，k % n 的操作利用了 同余 的概念。
	// 如果两个整数 a 和 b 除以 n 的余数相同，则称 a 和 b 对模 n 同余，记作 a ≡ b (mod n)。
	// 在旋转问题中，如果 k1 和 k2 对模 n 同余，则旋转 k1 位和旋转 k2 位得到的结果是相同的。
	// k = k % n 就是找到与原始 k 值对模 n 同余，且在 [-n+1, n-1] 范围内的值。
	k = k % n
	if k < 0 {
		k += n
	}
	fmt.Printf("rotate left n is %d, k is %d\n", n, k)
	// 1. 翻转 0 ～ k - 1的元素，【1，2，3，4，5】 => [2,1,3,4,5]
	reverse(s, 0, k-1)
	fmt.Println(s)
	// 2. 翻转 k ～ n - 1的元素，[2,1,3,4,5] => [2,1,5,4,3]
	reverse(s, k, n-1)
	fmt.Println(s)
	// 3. 翻转 0 ～ n - 1的元素，[2,1,5,4,3] => [3,4,5,1,2]
	reverse(s, 0, n-1)
}

// reverse 切片反转
func reverse(s []int, start, end int) {
	for i, j := start, end; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// 不满足就地移除
func DuplicateArray(s []string) []string {
	if len(s) == 0 || len(s) == 1 {
		return s
	}
	var ss = make([]string, 0)
	for i := 0; i < len(s); i++ {
		if i+1 >= len(s) {
			ss = append(ss, s[i])
			return ss
		}

		if s[i] != s[i+1] {
			ss = append(ss, s[i])
		}
	}

	return ss
}

// 移除重复元素（原地）
// 原理：双指针，i指向最后一个不重复元素的下标，j遍历当前元素；
// 所以，i=0，是第一个不重复的位置，j遍历数组，若 s[i] != s[j], 说明不重复，i 前移，当前j的位置是最后一个不重复的下标；
// 若 s[i] = s[j]，则重复，i 不移动，仍然指向 最后一个不重复的元素，j ++;
// 当 j 遍历完成后，截取 0 ～i的slice即可。
func RemoveAdjacentDuplicates(s []string) []string {
	if len(s) <= 1 {
		return s
	}
	//i标识应该被替代的位置
	i := 0
	for j := 1; j < len(s); j++ {
		if s[i] != s[j] {
			// 从0开始推演，i=0不需要被替代
			// i ++
			i++
			s[i] = s[j]
		}
	}
	return s[:i+1]
}

// reduceAdjacentSpaces 就地将 UTF-8 编码的字节 slice 中相邻的 Unicode 空白字符缩减为一个 ASCII 空白字符
func reduceAdjacentSpaces(b []byte) []byte {
	if len(b) == 0 {
		return b
	}

	i := 0 // 指向下一个非空白字符或 ASCII 空格应该放置的位置
	spaceFound := false

	for j := 0; j < len(b); {
		r, size := utf8.DecodeRune(b[j:])

		if unicode.IsSpace(r) {
			if !spaceFound {
				// 遇到第一个空白字符，将其替换为 ASCII 空格
				b[i] = ' '
				i++
				spaceFound = true
			}
		} else {
			// 遇到非空白字符，将其复制到 i 的位置
			for k := 0; k < size; k++ {
				b[i] = b[j+k]
				i++
			}
			spaceFound = false
		}
		j += size
	}

	return b[:i]
}

func Utf8Code() {
	// "你好，世界！" 的 UTF-8 编码
	utf8Bytes := []byte{
		0xE4, 0xBD, 0xA0, // 你 (U+4F60)
		0xE5, 0xA5, 0xBD, // 好 (U+597D)
		0xEF, 0xBC, 0x8C, // ， (U+FF0C)
		0xE4, 0xB8, 0x96, // 世 (U+4E16)
		0xE7, 0x95, 0x8C, // 界 (U+754C)
		0xEF, 0xBC, 0x81, // ！ (U+FF01)
	}

	fmt.Printf("字节 slice: % X\n", utf8Bytes) // 打印十六进制表示

	// 解码 UTF-8 字节 slice
	for i := 0; i < len(utf8Bytes); {
		r, size := utf8.DecodeRune(utf8Bytes[i:])
		fmt.Printf("字符: %c, 码点: U+%04X, 字节长度: %d\n", r, r, size)
		i += size
	}
	//将字节slice转为字符串
	str := string(utf8Bytes)
	fmt.Println("字符串", str)
}

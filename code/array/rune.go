package array

import "fmt"

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

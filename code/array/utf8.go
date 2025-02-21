package array

import (
	"fmt"
	"reflect"
	"unicode/utf8"
)

// Error
// 该版本实现存在问题，直接翻转[]byte后再使用utf8.DecodeRune会造成乱码，进而导致第二次翻转出错
func errorReverseUTF8StringInPlace(b []byte) []byte {
	// 1. 反转整个字节切片 (字节反转)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	// 2. 逐个 Rune 反转每个 Rune 的字节
	i := 0
	for i < len(b) {
		_, size := utf8.DecodeRune(b[i:]) // 解码当前位置的 Rune
		if size > 1 {
			// 如果 Rune 占用多个字节，反转这些字节
			reverseBytes(b[i : i+size])
		}
		i += size
	}
	return b
}

// 【修正后】，修正后的版本将[]byte直接转换为rune进行处理，避免错误
// reverseUTF8StringInPlace 就地反转 UTF-8 编码的字节 slice 中的字符（runes）
func ReverseUTF8StringInPlace(b []byte) []byte { //注意，这里返回了[]byte类型，但b的底层数据已经被修改
	// 1. 先解码成 runes
	runes := make([]rune, 0, len(b)) // 预估容量，避免多次扩容
	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		runes = append(runes, r)
		b = b[size:]
	}

	// 2. 反转 runes
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	// 3. 将反转后的 runes 编码回 UTF-8 字节
	result := make([]byte, 0, len(b)) // 初始容量与原始字节切片相同
	for _, r := range runes {
		buf := make([]byte, utf8.RuneLen(r)) // 为当前 rune 分配足够空间
		utf8.EncodeRune(buf, r)              // 将 rune 编码为 UTF-8 字节
		result = append(result, buf...)      // 追加到结果切片
	}

	return result
}

// reverseBytes 就地反转字节切片 (辅助函数)
func reverseBytes(b []byte) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}

//
//练习4.7：修改函数reverse，来翻转一个UTF-8编码的字符串中的字符元素，传入参数是该字符串对应的字节slice类型([]byte)。你可以做到不需要重新分配内存就实现该功能吗？
// 该问题需要注意的几个点：
// 1. 直接翻转非英文的字符可能会导致 utf8.DecodeRune 方法出错, 解决办法是先转为rune，再翻转
// 2. 原地翻转，

func UExamp() {
	testCases := map[string][]byte{
		"hello":      []byte("olleh"),
		"你好，世界":      []byte{0xE7, 0x95, 0x8C, 0xE4, 0xB8, 0x96, 0xEF, 0xBC, 0x8C, 0xE5, 0xA5, 0xBD, 0xE4, 0xBD, 0xA0},       // "界世，好你" 的 UTF-8 编码
		"Hello, 世界！": []byte{0xEF, 0xBC, 0x81, 0xE7, 0x95, 0x8C, 0xE4, 0xB8, 0x96, 0x20, 0x2C, 0x6F, 0x6C, 0x6C, 0x65, 0x48}, // "！界世 ,olleH" 的 UTF-8
		"ரேவதி":      []byte{0xBA, 0xB5, 0xBB, 0xBA, 0xBF, 0xAE},                                                             // "திவேரே" 的 UTF-8
		"😊👍😄":        []byte{0xF0, 0x9F, 0x98, 0x84, 0xF0, 0x9F, 0x91, 0x8D, 0xF0, 0x9F, 0x98, 0x8A},                         // "😄👍😊" 的 UTF-8
		"":           []byte{},
	}

	for tc, expected := range testCases {
		fmt.Printf("Original:   %q\n", tc)
		originalBytes := []byte(tc)
		reversedBytes := ReverseUTF8StringInPlace(originalBytes)
		fmt.Printf("Reversed: %q\n", string(reversedBytes))

		// 验证反转后的字节序列是否与预期相同
		if reflect.DeepEqual(reversedBytes, expected) {
			fmt.Println("Reversed bytes: OK")
		} else {
			fmt.Println("Reversed bytes: FAIL")
			fmt.Printf("  Expected: % X\n", expected)
			fmt.Printf("  Actual:   % X\n", reversedBytes)
		}

		fmt.Println()
	}
}

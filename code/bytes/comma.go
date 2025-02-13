package bytes

import (
	"bytes"
	"strings"
)

func Comma(s string) string {
	var buf bytes.Buffer
	n := len(s)
	if n == 0 {
		return ""
	}

	// 检查并处理符号
	start := 0
	if s[0] == '+' || s[0] == '-' {
		buf.WriteByte(s[0]) // 写入符号
		start = 1
	}

	// 计算第一个逗号前有多少位
	firstComma := (n - start) % 3
	if firstComma == 0 {
		firstComma = 3
	}

	// 写入第一个逗号前的数字
	buf.WriteString(s[start : start+firstComma])
	start += firstComma

	// 循环写入剩余的数字和逗号
	for start < n {
		buf.WriteByte(',')
		buf.WriteString(s[start : start+3])
		start += 3
	}

	return buf.String()
}

func NewComma(s string) string {
	var buf bytes.Buffer

	// 前置处理
	n := len(s)
	if n == 0 {
		return ""
	}

	start := 0
	end := n
	if s[0] == '-' || s[0] == '+' {
		buf.WriteByte(s[0])
		start = 1
	}

	//判断是否包含.
	dotIndex := strings.IndexByte(s, '.')
	if dotIndex != -1 {
		end = dotIndex
	}

	intPart := s[start:end]
	intLen := len(intPart)
	if intLen > 0 {
		// 获取第一个 ，位置
		firstComma := intLen % 3
		if firstComma == 0 {
			firstComma = 3
		}
		buf.WriteString(intPart[:firstComma])

		for i := firstComma; i < intLen; i += 3 {
			buf.WriteByte(',')
			buf.WriteString(intPart[i : i+3])
		}
	}

	if dotIndex != -1 {
		buf.WriteString(s[dotIndex:])
	}

	return buf.String()
}

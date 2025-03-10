package stringx

import "bytes"

func Comma(s string) string {
	n := len(s)
	if n == 0 {
		return s
	}
	var buf bytes.Buffer
	start := 0
	// , 是从右向左加的，所以我们进行取模运算
	firstComma := (n - start) % 3

	if firstComma == 0 {
		// 这里的意义是将长度恰好为3的情况排除，取模为0，
		firstComma = 3
	}
	// 追加字符串是从左向右添加
	// 先将 ，前的数字写入 ，
	buf.WriteString(s[start : start+firstComma])
	start += firstComma
	for start < n {
		buf.WriteByte(',')
		buf.WriteString(s[start : start+3])
		start += 3
	}

	return buf.String()
}

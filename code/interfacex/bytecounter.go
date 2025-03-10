package interfacex

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

type Counter struct {
	words, lines int
}

func (c *Counter) Write(p []byte) (int, error) {
	scanner := bufio.NewScanner(bytes.NewReader(p))

	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		c.words++
	}

	scanner = bufio.NewScanner(bytes.NewReader(p))
	for scanner.Scan() {
		c.lines++
	}

	return len(p), nil
}

func (c *Counter) Reset() {
	c.words = 0
	c.lines = 0
}

func TestCounter() {
	var c = new(Counter)

	c.Write([]byte("\n"))

	fmt.Println(c)

	c.Reset()

	c.Write([]byte("hello world\n ni hao shi jie"))
	fmt.Println(c)
}

type ByteCounter int

func (c *ByteCounter) Write(p []byte) (int, error) {
	*c += ByteCounter(len(p))
	return len(p), nil
}

func TestByteCounter() {
	var c ByteCounter
	c.Write([]byte("hello"))

	fmt.Println(c)

	c = 0
	var name = "Dolly"
	fmt.Fprintf(&c, "hello, %s", name)
	fmt.Println(c)
}

// 实现一个满足如下签名的CountingWriter函数，
// 输入一个io.Writer，输出一个封装了输入值的新Writer，以及一个指向int64的指针，该指针对应的值是新的Writer写入的字节数。
// 函数签名：func CounterWriter(w io.Writer) (io.Writer,*int64)
func CountingWriter(w io.Writer) (io.Writer, *int64) {
	var count int64
	cw := &countingWriter{w: w, count: &count}
	return cw, &count
}

type countingWriter struct {
	w     io.Writer
	count *int64
}

func (cw *countingWriter) Write(p []byte) (int, error) {
	n, err := cw.w.Write(p)
	*cw.count += int64(n) // 将写入的字节数累加到计数器
	return n, err
}

func TestCountingWriter() {
	// 示例 1：写入到 os.Stdout
	writer, count := CountingWriter(os.Stdout)
	fmt.Fprintf(writer, "Hello, world!\n")
	fmt.Println("Bytes written:", *count) // 输出写入的字节数

	// 多次写入
	writer, count = CountingWriter(os.Stdout)
	fmt.Fprintf(writer, "First line\n")
	fmt.Fprintf(writer, "Second line\n")
	fmt.Println("Total bytes written:", *count)
}

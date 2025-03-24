package interfacex

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// 接口赋值
func TestInterface() {
	var w io.Writer
	w = os.Stdout
	w = new(bytes.Buffer)
	w = nil

	fmt.Println(w)
}

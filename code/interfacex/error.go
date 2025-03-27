package interfacex

import (
	"errors"
	"fmt"
	"syscall"
)

func MyError() {
	errors.New("EOF")
}

func SysError() {
	ss := [...]string{
		1: "hello",
		2: "nihao",
	}

	fmt.Println(ss[0], ss[1], ss[2])
}

func SysError2() {
	var err error = syscall.Errno(2)
	fmt.Println(err.Error())
	fmt.Println(err)
}

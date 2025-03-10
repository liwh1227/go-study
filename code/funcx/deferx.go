package funcx

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

func fetch2(url string) (string, int64, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}

	defer resp.Body.Close()

	local := path.Base(resp.Request.URL.Path)
	if local == "/" {
		local = "index.html"
	}

	f, err := os.Create(local)
	if err != nil {
		return "", 0, err
	}

	n, err := io.Copy(f, resp.Body)
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}

	return local, n, err
}

func fetch(url string) (string, int64, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			fmt.Printf("record err %v", err)
			return
		}
	}()

	local := path.Base(resp.Request.URL.Path)
	if local == "/" {
		local = "index.html"
	}

	f, err := os.Create(local)
	if err != nil {
		return "", 0, err
	}

	n, err := io.Copy(f, resp.Body)
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}

	return local, n, err
}

func DoDefer() {
	outer()
}

func DoDefer2() {
	defer func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Inner recover:", r)
			}
		}()

		panic("second panic")
	}()

	panic("first panic")
}

func outer() {
	defer fmt.Println("Outer 1") // 属于 outer 的堆栈, 第 3 步

	inner := func() {
		defer fmt.Println("Inner 1") // 属于 inner 的堆栈, 第 1 步
		defer fmt.Println("Inner 2") // 属于 inner 的堆栈, 第 2 步
	}
	inner()

	defer fmt.Println("Outer 2") // 属于 outer 的堆栈, 第 4 步
}

func pr() {
	p := recover()
	switch p {
	case nil:

	}

}

package main

import (
	"fmt"
	"github.com/pkg/errors"
)

func main() {
	var res = make([]interface{}, 0)

	res = append(res, 0, 1, 23, 4)

	for i := range res {
		fmt.Println(res[i])
	}
}

func getNum(i int) (int, error) {
	if i == 4 {
		return i, errors.New("get num has error")
	}

	return i, nil
}

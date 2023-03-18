package main

import (
	"errors"
	"fmt"
	"src/chainmaker-client"
	"sync"
)

//

func main() {
	wg := sync.WaitGroup{}

	wg.Add(5)

	go func() {
		fmt.Println("go 1")
		defer wg.Done()
		_, err := chainmaker_client.NewClient()
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	go func() {
		fmt.Println("go 2")
		defer wg.Done()
		_, err := chainmaker_client.NewClient()
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	go func() {
		fmt.Println("go 3")
		defer wg.Done()
		_, err := chainmaker_client.NewClient()
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	go func() {
		fmt.Println("go 4")
		defer wg.Done()
		_, err := chainmaker_client.NewClient()
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	go func() {
		fmt.Println("go 5")
		defer wg.Done()
		_, err := chainmaker_client.NewClient()
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	wg.Wait()

	fmt.Println("xxxxxx end xxxxxx")
}

func func1() {
	fmt.Println("i am func 1")
	func2()
}

func func2() {
	fmt.Println("i am func 2")
	func3()
}

func func3() {
	fmt.Println("i am func 3")
	panic(errors.New("unknown error, panic"))
}

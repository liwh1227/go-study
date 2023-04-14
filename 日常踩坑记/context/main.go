package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer wg.Done()
		fmt.Println("wait cancel....")
		<-ctx.Done()

		fmt.Println("i am children, i am done ")
		return
	}()

	fmt.Println("start cancel ")
	time.Sleep(5 * time.Second)
	cancel()

	wg.Wait()

	fmt.Println("i am done...")
}

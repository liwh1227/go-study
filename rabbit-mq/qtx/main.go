package main

import (
	"fmt"
	"qtx/rabbitmq"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	fmt.Println("start api service success.")

	conf := rabbitmq.Conf{
		User:  "rabbit",
		Pwd:   "rabbit@123",
		Addr:  "127.0.0.1",
		Port:  "5672",
		Vhost: "gateway-dev",
	}

	err := rabbitmq.Init(conf)
	if err != nil {
		fmt.Println(err)
		return
	}

	queueName := "gateway.user.queue"

	wg.Add(1)
	go func() {
		if err := rabbitmq.NewConsumer(queueName, func(body []byte) error {
			//if string(body) == "success" {
			fmt.Println("consume msg :" + string(body))
			return nil
			//}
			//return errors.New("consume msg failed!")
		}); err != nil {
			wg.Done()
			fmt.Println(err)
		}
	}()

	wg.Wait()

	fmt.Println("Main process exit")
}

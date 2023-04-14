package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"qtx/api"
	"qtx/rabbitmq"
	"syscall"
	"time"
)

func main() {
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

	// 创建兑换的队列
	err = rabbitmq.CreateExchangeQueue()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 创建气泡队列
	err = rabbitmq.CreateBubbleQueue()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	var errChan chan error
	go api.Run(ctx, errChan)

	// 主动退出信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		cancel()
	case herr := <-errChan:
		fmt.Println(herr)

		time.Now().Unix()
	}

	fmt.Println("Main process exit")
}

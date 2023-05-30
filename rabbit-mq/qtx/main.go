package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/signal"
	"qtx/api"
	"qtx/rabbitmq"
	"sync"
	"syscall"
)

const (
	q1 = "qtx.bc.exchange.queue"
	q2 = "qtx.bc.bubble.queue"
)

var conf = rabbitmq.Conf{
	User:  "rabbit",
	Pwd:   "rabbit@123",
	Addr:  "127.0.0.1",
	Port:  "5672",
	Vhost: "gateway-dev",
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := consumer("consumer1", handleExchange)
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	fmt.Println("start success...")
	wg.Wait()
}

func producer() {
	fmt.Println("start api service success.")
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
		cancel()
	}
}

func consumer(name string, cb func(msg []byte) error) error {
	m, err := rabbitmq.NewMQ(rabbitUrl(), conf.Vhost).Open()
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer m.Close()

	c, err := m.Consumer(name)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer c.Close()

	msgC := make(chan rabbitmq.Delivery, 1)
	defer close(msgC)

	// 绑定消费队列
	err = c.SetQueueBinds([]string{q1, q2}).SetMsgCallback(msgC).OpenWithDirectConsumeType()
	if err != nil {
		fmt.Println(err)
		return err
	}

	for msg := range msgC {
		err = cb(msg.Body)
		if err != nil {
			_ = msg.Reject(false)
			fmt.Println(err)
			continue
		}

		err = msg.Ack(false)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func rabbitUrl() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", conf.User, conf.Pwd, conf.Addr, conf.Port)
}

// 处理兑换的消费消息
func handleExchange(msg []byte) error {
	fmt.Printf("[consumer] start consume msg %s\n", string(msg))
	return errors.New("handle error")
}

func handleBubble(msg []byte) error {
	fmt.Printf("<=====[consumer] start consume msg %s\n", string(msg))
	return nil
}

package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	BubbleExchange = "qtx.bc.bubble.exchange"
	BubbleQueue    = "bubble_chain"
	BubbleKey      = "qtx.bc.bubble.key"

	ExchangeExchange = "update_exchange_status"
	ExchangeQueue    = "update_exchange_status"
	ExchangeKey      = "qtx.bc.exchange.key"
)

// 创建气泡的队列
func CreateBubbleQueue() error {
	// 创建队列，并绑定对应死信的key和exhcange
	err := NewChannel().QueueDeclare(BubbleQueue)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = NewChannel().ExchangeDeclare(BubbleExchange, amqp.ExchangeDirect)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = NewChannel().QueueBind(BubbleQueue, BubbleKey, BubbleExchange)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// 创建兑换的队列
func CreateExchangeQueue() error {
	// 创建队列
	err := NewChannel().QueueDeclare(ExchangeQueue)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = NewChannel().ExchangeDeclare(ExchangeExchange, amqp.ExchangeDirect)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = NewChannel().QueueBind(ExchangeQueue, ExchangeKey, ExchangeExchange)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

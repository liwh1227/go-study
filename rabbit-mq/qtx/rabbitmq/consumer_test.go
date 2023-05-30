package rabbitmq

import (
	"fmt"
	"github.com/pkg/errors"
	"testing"
	"time"
)

var (
	LocalMqURL = "amqp://rabbit:rabbit@123@127.0.0.1:5672/"
	QtxMqUrl   = "amqp://htxtest:123456@123@192.168.0.72:5672/"
)

func TestConsumer(t *testing.T) {

	m, err := NewMQ(LocalMqURL, "gateway-dev").Open()
	if err != nil {
		t.Error(err)
		return
	}

	defer m.Close()

	c, err := m.Consumer("test-consume")
	if err != nil {
		panic(fmt.Sprintf("Create consumer failed, %v", err))
	}
	defer c.Close()

	//exb := []*ExchangeBinds{
	//	&ExchangeBinds{
	//		Exch: DefaultExchange("", amqp.ExchangeDirect),
	//		Bindings: []*Binding{
	//			&Binding{
	//				Queues: []*Queue{
	//					DefaultQueue("bubble_chain"),
	//				},
	//			},
	//			//&Binding{
	//			//	Queues: []*Queue{
	//			//		DefaultQueue("qtx.bc.exchange.queue"),
	//			//	},
	//			//},
	//		},
	//	},
	//}
	msgC := make(chan Delivery, 1)
	defer close(msgC)

	err = c.SetQueueBinds([]string{
		"qtx.bc.exchange.queue",
		"qtx.bc.bubble.queue",
	}).SetMsgCallback(msgC).Open()
	if err != nil {
		t.Error(err)
		return
	}

	for msg := range msgC {
		t.Logf("Consumer receive msg `%s`，from queue name is %v,exchange is %v\n", string(msg.Body), msg.RoutingKey, msg.Exchange)
		time.Sleep(1 * time.Second)

		err = msg.Ack(false)
		if err != nil {
			t.Error(err)
			continue
		}
	}
}

func TestDirectConsumer(t *testing.T) {
	m, err := NewMQ(LocalMqURL, "gateway-dev").Open()
	if err != nil {
		t.Error(err)
		return
	}

	defer m.Close()

	c, err := m.Consumer("")
	if err != nil {
		t.Error(err)
		return
	}

	defer c.Close()
	msgC := make(chan Delivery, 1)
	defer close(msgC)

	var (
		q1 = "qtx.bc.exchange.queue"
		q2 = "qtx.bc.bubble.queue"
	)

	err = c.SetQueueBinds([]string{q1, q2}).SetMsgCallback(msgC).OpenWithDirectConsumeType()
	if err != nil {
		t.Error(err)
		return
	}

	params := make([]*RegisterHandlerParam, 0)
	param1 := &RegisterHandlerParam{
		QueueName: q1,
		F:         handleExchange,
	}
	param2 := &RegisterHandlerParam{
		QueueName: q2,
		F:         handleBubble,
	}

	params = append(params, param1, param2)

	err = c.RegisterHandler(params...)
	if err != nil {
		t.Error(err)
		return
	}

	for msg := range msgC {
		fmt.Println(msg.RoutingKey)
		handler := HandleFuncMap[msg.RoutingKey]
		if handler != nil {
			err = handler(msg.Body)
			if err != nil {
				_ = msg.Reject(false)
				t.Error(err)
				continue
			}

			err = msg.Ack(false)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("pls register msg handle func...")
		}

	}
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

package rabbitmq

import (
	"fmt"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

const (
	stateClosed = iota
	stateOpened
	stateReopening
)

// 通过queue增加consumer的处理函数
func NewConsumerWithQueue(queue string, handler func(msg []byte) error) error {
	c, err := GetMq().Consumer(queue)
	if err != nil {
		return err
	}

	msgC := make(chan Delivery, 1)
	err = c.setQueueBinds(queue).setMsgCallback(msgC).open()
	if err != nil {
		return err
	}

	for msg := range msgC {
		err = handler(msg.Body)
		if err != nil {
			err = msg.Reject(false)
			if err != nil {
				log.Println(err)
			}
			continue
		}
		err = msg.Ack(false)
		if err != nil {
			log.Println(err)
		}
	}

	return err
}

// 通过声明的方式创建消费者，如果队列和excahnge不存在，会直接创建
func NewConsumerWithDeclare(exchangeBinds []*ExchangeBinds, handler func(msg []byte) error) error {
	c, err := GetMq().Consumer("")
	if err != nil {
		return err
	}

	msgC := make(chan Delivery, 1)
	err = c.setExchangeBinds(exchangeBinds).setMsgCallback(msgC).open()
	if err != nil {
		return err
	}

	for msg := range msgC {
		err = handler(msg.Body)
		if err != nil {
			err = msg.Reject(false)
			if err != nil {
				log.Println(err)
			}
			continue
		}
		err = msg.Ack(false)
		if err != nil {
			log.Println(err)
		}
	}

	return err
}

// ExchangeBinds 绑定顺序：exchange ==> routeKey ==> queues
type ExchangeBinds struct {
	Exch     *Exchange
	Bindings []*Binding
}

// Exchange 基于amqp的Exchange配置
type Exchange struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table // default is nil
}

// Biding routeKey ==> queues
type Binding struct {
	RouteKey string
	Queues   []*Queue
	NoWait   bool       // default is false
	Args     amqp.Table // default is nil
}

// Queue 基于amqp的Queue配置
type Queue struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

// 创建默认的exchange
func DefaultExchange(name string, kind string) *Exchange {
	return &Exchange{
		Name:       name,
		Kind:       kind,
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
		Args:       nil,
	}
}

// 创建默认队列
func DefaultQueue(name string) *Queue {
	return &Queue{
		Name:       name,
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
	}
}

// 死信队列
func DeadLetterQueue(name, exchange, routingKey string) *Queue {
	return &Queue{
		Name:       name,
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args: amqp.Table{
			"x-dead-letter-exchange":    exchange,
			"x-dead-letter-routing-key": routingKey,
		},
	}
}

// ch绑定exchange，建立exchange、routing key及queue之间的关系
func applyExchangeBinds(ch *amqp.Channel, exchangeBinds []*ExchangeBinds) (err error) {
	if ch == nil {
		return errors.New("MQ: Nil producer channel")
	}
	if len(exchangeBinds) <= 0 {
		return errors.New("MQ: Empty exchangeBinds")
	}

	for _, eb := range exchangeBinds {
		if eb.Exch == nil {
			return errors.New("MQ: Nil exchange found.")
		}
		if len(eb.Bindings) <= 0 {
			return fmt.Errorf("MQ: No bindings queue found for exchange(%s)", eb.Exch.Name)
		}
		// declare exchange
		if err = ch.ExchangeDeclare(eb.Exch.Name, eb.Exch.Kind, eb.Exch.Durable, eb.Exch.AutoDelete, eb.Exch.Internal, eb.Exch.NoWait, eb.Exch.Args); err != nil {
			return fmt.Errorf("MQ: Declare exchange(%s) failed, %v", eb.Exch.Name, err)
		}

		// declare and bind queues
		for _, b := range eb.Bindings {
			if b == nil {
				return fmt.Errorf("MQ: Nil binding found, exchange:%s", eb.Exch.Name)
			}
			if len(b.Queues) <= 0 {
				return fmt.Errorf("MQ: No queues found for exchange(%s)", eb.Exch.Name)
			}
			for _, q := range b.Queues {
				if q == nil {
					return fmt.Errorf("MQ: Nil queue found, exchange:%s", eb.Exch.Name)
				}
				if _, err = ch.QueueDeclare(q.Name, q.Durable, q.AutoDelete, q.Exclusive, q.NoWait, q.Args); err != nil {
					return fmt.Errorf("MQ: Declare queue(%s) failed, %v", q.Name, err)
				}
				// 直接绑定队列进行消费
				if err = ch.QueueBind(q.Name, b.RouteKey, eb.Exch.Name, b.NoWait, b.Args); err != nil {
					return fmt.Errorf("MQ: Bind exchange(%s) <--> queue(%s) failed, %v", eb.Exch.Name, q.Name, err)
				}
			}
		}
	}
	return nil
}

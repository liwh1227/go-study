package rabbitmq

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type (
	Conf struct {
		Addr       string
		Port       string
		User       string
		Pwd        string
		Vhost      string
		MaxRetry   int
		RetryTimes int
	}
)

// Connection amqp.Connection wrapper
type Connection struct {
	*amqp.Connection
}

func (c *Connection) Channel() (*Channel, error) {
	ch, err := c.Connection.Channel()
	if err != nil {
		return nil, err
	}

	channel := &Channel{
		Channel: ch,
	}

	go func() {
		for {
			fmt.Println("will listen amqp error msg...")
			reason, ok := <-channel.Channel.NotifyClose(make(chan *amqp.Error))
			fmt.Println("wait recreate channel...")
			// exit this goroutine if closed by developer
			if !ok || channel.IsClosed() {
				log.Println("channel closed")
				_ = channel.Close() // close again, ensure closed flag set when connection closed
				break
			}
			log.Printf("channel closed, reason: %v", reason)

			// reconnect if not closed by developer
			for {
				// wait 1s for connection reconnect
				time.Sleep(3 * time.Second)

				ch, err := c.Connection.Channel()
				if err == nil {
					log.Println("channel recreate success")
					channel.Channel = ch
					break
				}

				log.Printf("channel recreate failed, err: %v", err)
			}
		}

	}()

	return channel, nil
}

type Channel struct {
	*amqp.Channel
	closed int32
}

// 默认rabbitmq连接
var (
	defaultConn    *Connection
	defaultChannel *Channel
)

func dial(url, vhost string) (*Connection, error) {
	conn, err := amqp.DialConfig(url, amqp.Config{
		Vhost: vhost,
	})
	if err != nil {
		return nil, err
	}

	fmt.Println(conn.Config.Vhost)

	connection := &Connection{
		Connection: conn,
	}

	go func() {
		for {
			reason, ok := <-connection.Connection.NotifyClose(make(chan *amqp.Error))
			// exit this goroutine if closed by developer
			if !ok {
				log.Println("connection closed")
				break
			}
			log.Printf("connection closed, reason: %v", reason)

			// reconnect if not closed by developer
			for {
				// wait 1s for reconnect
				time.Sleep(3 * time.Second)

				conn, err := amqp.DialConfig(url, amqp.Config{
					Vhost: vhost,
				})
				if err == nil {
					connection.Connection = conn
					log.Println("reconnect success")
					break
				}

				log.Printf("reconnect failed, err: %v", err)
			}
		}
	}()

	return connection, nil
}

// Init 初始化
func Init(c Conf) (err error) {
	if c.Addr == "" {
		return nil
	}
	defaultConn, err = dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
		c.User,
		c.Pwd,
		c.Addr,
		c.Port), c.Vhost)
	if err != nil {
		return fmt.Errorf("new mq conn err: %v", err)
	}

	defaultChannel, err = defaultConn.Channel()
	if err != nil {
		return fmt.Errorf("new mq channel err: %v", err)
	}

	return
}

func NewChannel() *Channel {
	return defaultChannel
}

// ExchangeDeclare 创建交换机.
func (ch *Channel) ExchangeDeclare(name string, kind string) (err error) {
	return ch.Channel.ExchangeDeclare(name, kind, true, false, false, false, nil)
}

// Publish 发布消息.
func (ch *Channel) Publish(exchange, key string, body []byte) (err error) {
	_, err = ch.Channel.PublishWithDeferredConfirmWithContext(context.Background(), exchange, key, false, false,
		amqp.Publishing{ContentType: "text/plain", Body: body})
	return err
}

// QueueBind 绑定队列.
func (ch *Channel) QueueBind(name, key, exchange string) (err error) {
	return ch.Channel.QueueBind(name, key, exchange, false, nil)
}

// QueueDeclareWithDelay 创建延迟队列
// ps: 如果exchange、key为其他队列的，那么这里创建的就是name的死信队列
func (ch *Channel) QueueDeclareWithDelay(name, exchange, key string) (err error) {
	_, err = ch.Channel.QueueDeclare(name, true, false, false, false, amqp.Table{
		"x-dead-letter-exchange":    exchange,
		"x-dead-letter-routing-key": key,
	})
	return
}

// QueueDeclare 创建队列.
func (ch *Channel) QueueDeclare(name string) (err error) {
	_, err = ch.Channel.QueueDeclare(name, true, false, false, false, nil)
	return
}

// NewConsumer 实例化一个消费者, 会单独用一个channel.
func NewConsumer(queue string, handler func([]byte) error) error {
	ch, err := defaultConn.Channel()
	if err != nil {
		return fmt.Errorf("new mq channel err: %v", err)
	}

	deliveries, err := ch.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume err: %v, queue: %s", err, queue)
	}

	for msg := range deliveries {
		err = handler(msg.Body)
		if err != nil {
			_ = msg.Reject(true)
			continue
		}
		_ = msg.Ack(false)
	}

	return nil
}

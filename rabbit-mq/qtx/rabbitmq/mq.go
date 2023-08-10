package rabbitmq

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var mqIns *MQ

var ()

type MQ struct {
	url       string           // mq 链接的url地址
	vhost     string           // vhost 名称
	mutex     sync.RWMutex     // 读写锁
	conn      *amqp.Connection // mq链接
	consumers []*Consumer      // mq:consumer 1 : N
	closeC    chan *amqp.Error // 捕捉链接错误
	stopC     chan struct{}    // 关闭通道
	state     int              // MQ状态
}

// 消息队列的封装
func InitRabbitMq() (err error) {
	mqIns, err = newMQ(fmt.Sprintf("amqp://%s:%s@%s:%s/",
		conf.User,
		conf.Pwd,
		conf.Addr,
		conf.Port), conf.Vhost).open()
	if err != nil {
		return fmt.Errorf("new mq conn err: %v", err)
	}

	return
}

func GetMq() *MQ {
	return mqIns
}

func newMQ(url, vhost string) *MQ {
	return &MQ{
		url:   url,
		vhost: vhost,
		state: stateClosed,
	}
}

// 建立和打开mq链接
func (m *MQ) open() (*MQ, error) {
	var err error
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 处于打开的状态
	if m.state == stateOpened {
		return m, errors.New("MQ: had been opened")
	}

	m.conn, err = m.dial()
	if err != nil {
		return m, fmt.Errorf("MQ: Dial err: %v", err)
	}

	m.state = stateOpened
	m.stopC = make(chan struct{})
	m.closeC = make(chan *amqp.Error, 1)
	m.conn.NotifyClose(m.closeC)

	// 处理链接中的错误情况
	go m.keepalive()

	return m, err
}

// 通过mq获取消费者
func (m *MQ) Consumer(name string) (*Consumer, error) {
	m.mutex.Lock()
	m.mutex.Unlock()

	// 查看当前队列的状态
	if m.state != stateOpened {
		return nil, errors.New("mq is not open state")
	}

	c := newConsumer(name, m)
	m.consumers = append(m.consumers, c)
	return c, nil
}

func (m *MQ) Close() {
	m.mutex.Lock()

	// close consumers
	for _, c := range m.consumers {
		c.close()
	}
	m.consumers = m.consumers[:0]

	// close mq connection
	select {
	case <-m.stopC:
		// had been closed
	default:
		close(m.stopC)
	}

	m.mutex.Unlock()

	// wait done
	for m.State() != stateClosed {
		time.Sleep(time.Second)
	}
}

func (m *MQ) State() int {
	return m.state
}

func (m *MQ) keepalive() {
	select {
	case <-m.stopC:
		// 正常关闭
		log.Println("MQ: Shutdown normally.")
		m.mutex.Lock()
		m.conn.Close()
		m.state = stateClosed
		m.mutex.Unlock()

	case err := <-m.closeC:
		if err == nil {
			log.Println("MQ: Disconnected with MQ, but Error detail is nil")
		} else {
			log.Printf("MQ: Disconnected with MQ, code:%d, reason:%s\n", err.Code, err.Reason)
		}

		// tcp连接中断, 重新连接
		m.mutex.Lock()
		m.state = stateReopening
		m.mutex.Unlock()

		maxRetry := 100
		for i := 0; i < maxRetry; i++ {
			time.Sleep(3 * time.Second)
			if _, e := m.open(); e != nil {
				log.Printf("MQ: Connection recover failed for %d times, %v\n", i+1, e)
				continue
			}
			log.Printf("MQ: Connection recover OK. Total try %d times\n", i+1)
			return
		}
		log.Printf("MQ: Try to reconnect to MQ failed over maxRetry(%d), so exit.\n", maxRetry)
	}
}

// 获取当前的链接的channel
func (m *MQ) channel() (*amqp.Channel, error) {
	return m.conn.Channel()
}

// 同mq server建立链接
func (m *MQ) dial() (*amqp.Connection, error) {
	return amqp.DialConfig(m.url, amqp.Config{Vhost: m.vhost})
}

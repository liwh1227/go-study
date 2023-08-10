package rabbitmq

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

var conf = Conf{
	User:       "rabbit",
	Pwd:        "rabbit@123",
	Addr:       "127.0.0.1",
	Port:       "5672",
	Vhost:      "gateway-dev",
	MaxRetry:   100,
	RetryTimes: 10,
}

type Delivery struct {
	amqp.Delivery
}

// 消费者消费选项
type consumeOption struct {
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}

// 默认的消费者参数
func defaultConsumeOption() *consumeOption {
	return &consumeOption{
		NoWait: true,
	}
}

// 消费者
type Consumer struct {
	name          string           // Consumer的名字, "" is OK
	mq            *MQ              // MQ实例
	mutex         sync.RWMutex     // 保护数据并发安全
	ch            *amqp.Channel    // MQ的会话channel
	exchangeBinds []*ExchangeBinds // MQ的exchange与其绑定的queues
	queueList     []string         // 直接根据队列名称进行消费
	prefetch      int              // Ops prefetch
	callback      chan<- Delivery  // 上层用于接收消费出来的消息的管道
	closeC        chan *amqp.Error // 监听会话channel关闭，mq意外退出的情况
	stopC         chan struct{}    // Consumer关闭控制，主动关闭mq
	state         int              // Consumer状态
	queue         string           // 队列名称
}

func newConsumer(name string, mq *MQ) *Consumer {
	return &Consumer{
		name:   name,
		mq:     mq,
		stopC:  make(chan struct{}),
		closeC: make(chan *amqp.Error, 1),
	}
}

func (c *Consumer) Name() string {
	return c.name
}

func (c *Consumer) closeChan() {
	c.mutex.Lock()
	c.ch.Close()
	c.mutex.Unlock()
}

// 具体接受的管道信息
func (c *Consumer) setQueueBinds(queue string) *Consumer {
	c.mutex.Lock()
	c.queue = queue
	c.mutex.Unlock()
	return c
}

// 设置
func (c *Consumer) setExchangeBinds(eb []*ExchangeBinds) *Consumer {
	c.mutex.Lock()
	if c.state != stateOpened {
		c.exchangeBinds = eb
	}
	c.mutex.Unlock()
	return c
}

// 具体接受的管道信息
func (c *Consumer) setMsgCallback(cb chan<- Delivery) *Consumer {
	c.mutex.Lock()
	c.callback = cb
	c.mutex.Unlock()
	return c
}

// 设置qos
func (c *Consumer) setQos(prefetch int) *Consumer {
	c.mutex.Lock()
	c.prefetch = prefetch
	c.mutex.Unlock()
	return c
}

// 该模式下直接针对队列名称进行消费，无须绑定exchange
func (c *Consumer) open() error {
	// Open期间不允许对channel做任何操作
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 参数校验
	if c.mq == nil {
		return errors.New("MQ: Bad consumer,can not find any mq")
	}

	// 状态检测
	if c.state == stateOpened {
		return errors.New("MQ: Consumer had been opened")
	}

	// 初始化该consumer的channel，每个消费者独占一个channel
	ch, err := c.mq.channel()
	if err != nil {
		return errors.Errorf("MQ: create channel failed %v", err)
	}

	c.ch = ch
	c.state = stateOpened
	c.ch.NotifyClose(c.closeC)

	opt := defaultConsumeOption()
	notify := make(chan error, 1)

	c.consume(opt, notify)
	for e := range notify {
		if e != nil {
			fmt.Printf("[ERROR] consume has some err %v\n", e)
			continue
		}
		break
	}
	close(notify)

	go c.keepalive(c.open)

	return nil
}

// 消费者主动关闭
func (c *Consumer) close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	select {
	case <-c.stopC:
		// had been closed
	default:
		close(c.stopC)
	}
}

// 处理接受到的队列信息
func (c *Consumer) deliver(delivery <-chan amqp.Delivery) {
	for d := range delivery {
		if c.callback != nil {
			c.callback <- Delivery{d}
		}
	}
}

// 获取当前消费者状态
func (c *Consumer) State() int {
	return c.state
}

// 监听当前mq的状态，并进行重连
func (c *Consumer) keepalive(reopen func() error) {
	select {
	case <-c.stopC:
		c.mutex.Lock()
		c.ch.Close()
		c.ch = nil
		c.state = stateClosed
		c.mutex.Unlock()

	case err := <-c.closeC:
		if err == nil {
			log.Printf("MQ: Consumer(%s)'s channel was closed, but Error detail is nil\n", c.name)
		} else {
			log.Printf("MQ: Consumer(%s)'s channel was closed, code:%d, reason:%s\n", c.name, err.Code, err.Reason)
		}

		// channel被异常关闭了
		c.mutex.Lock()
		c.state = stateReopening
		c.mutex.Unlock()

		for i := 0; i < conf.MaxRetry; i++ {
			time.Sleep(time.Duration(conf.RetryTimes) * time.Second)
			// 1. 判断mq是否处于开启状态
			if c.mq.State() != stateOpened {
				log.Printf("MQ: Consumer(%s) try to recover channel for %d times, but mq's state != stateOpened\n", c.name, i+1)
				continue
			}

			e := reopen()
			if e != nil {
				log.Printf("MQ: Consumer(%s) recover channel failed for %d times, Err:%v\n", c.name, i+1, e)
				continue
			}

			// 3. 重新打开成功
			log.Printf("MQ: Consumer(%s) recover channel OK. Total try %d times\n", c.name, i+1)
			return
		}
		// 超过最大重启次数
		log.Printf("MQ: Consumer(%s) try to recover channel over maxRetry(%d), so exit\n", c.name, conf.MaxRetry)
	}
}

// 该方法需要进行exchange routing的绑定
func (c *Consumer) openWithDeclare() error {
	// Open期间不允许对channel做任何操作
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 参数校验
	if c.mq == nil {
		return errors.New("MQ: Bad consumer,can not find any mq")
	}

	// 未进行exchange的绑定
	if len(c.exchangeBinds) <= 0 {
		return errors.New("MQ: No exchangeBinds found. You should SetExchangeBinds brefore open.")
	}

	// 状态检测
	if c.state == stateOpened {
		return errors.New("MQ: Consumer had been opened")
	}

	// 初始化该consumer的channel，每个消费者独占一个channel
	ch, err := c.mq.channel()
	if err != nil {
		return errors.Errorf("create channel failed %v", err)
	}

	// 绑定exchange
	err = func(ch *amqp.Channel) error {
		var e error
		if e = applyExchangeBinds(ch, c.exchangeBinds); e != nil {
			return e
		}
		if e = ch.Qos(c.prefetch, 0, false); e != nil {
			return e
		}
		return nil
	}(ch)
	if err != nil {
		return errors.Errorf("MQ: %v", err)
	}

	c.ch = ch
	c.state = stateOpened
	c.ch.NotifyClose(c.closeC)

	opt := defaultConsumeOption()
	notify := make(chan error, 1)

	c.consume(opt, notify)
	for e := range notify {
		if e != nil {
			log.Printf("consume has some err %v\n", e)
			continue
		}
		break
	}
	close(notify)

	go c.keepalive(c.openWithDeclare)

	return nil
}

// notifyErr 向上层抛出错误, 如果error为空表示执行完成.由上层负责关闭channel
func (c *Consumer) consume(opt *consumeOption, notifyErr chan<- error) {
	if len(c.queue) != 0 {
		// 说明直接绑定某队列进行消费，无须声明
		delivery, err := c.ch.Consume(c.queue, "", opt.AutoAck, opt.Exclusive, opt.NoLocal, opt.NoWait, opt.Args)
		if err != nil {
			notifyErr <- fmt.Errorf("MQ: Consumer(%s) consume queue(%s) failed, %v", c.name, c.queue, err)
		}
		go c.deliver(delivery)
	} else {
		for idx, eb := range c.exchangeBinds {
			if eb == nil {
				notifyErr <- fmt.Errorf("MQ: ExchangeBinds[%d] is nil, consumer(%s)", idx, c.name)
				continue
			}
			for i, b := range eb.Bindings {
				if b == nil {
					notifyErr <- fmt.Errorf("MQ: Binding[%d] is nil, ExchangeBinds[%d], consumer(%s)", i, idx, c.name)
					continue
				}
				for qi, q := range b.Queues {
					if q == nil {
						notifyErr <- fmt.Errorf("MQ: Queue[%d] is nil, ExchangeBinds[%d], Biding[%d], consumer(%s)", qi, idx, i, c.name)
						continue
					}
					// 监听消费者绑定的队列，获取通道中的消费信息
					delivery, err := c.ch.Consume(q.Name, "", opt.AutoAck, opt.Exclusive, opt.NoLocal, opt.NoWait, opt.Args)
					if err != nil {
						notifyErr <- fmt.Errorf("MQ: Consumer(%s) consume queue(%s) failed, %v", c.name, q.Name, err)
						continue
					}
					go c.deliver(delivery)
				}
			}
		}
	}

	notifyErr <- nil
}

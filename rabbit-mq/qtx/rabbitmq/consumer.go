package rabbitmq

import (
	"fmt"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
	"time"
)

type msgHandler func(msg []byte) error

var HandleFuncMap = make(map[string]msgHandler, 0) // 该消费者需要处理的函数方法

type Delivery struct {
	amqp.Delivery
}

type RegisterHandlerParam struct {
	ConsumerName string
	Cb           msgHandler
}

// 消费者消费选项
type ConsumeOption struct {
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}

// 默认的消费者参数
func DefaultConsumeOption() *ConsumeOption {
	return &ConsumeOption{
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
}

func newConsumer(name string, mq *MQ) *Consumer {
	return &Consumer{
		name:  name,
		mq:    mq,
		stopC: make(chan struct{}),
	}
}

// 注册consumer要处理的消息方法，该方法与队列绑定
func (c *Consumer) registerHandler(params ...*RegisterHandlerParam) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for i := range params {
		_, ok := HandleFuncMap[params[i].ConsumerName]
		if ok {
			return errors.New("duplicated register handler")
		}
		HandleFuncMap[params[i].ConsumerName] = params[i].Cb
	}

	return nil
}

func (c *Consumer) Name() string {
	return c.name
}

func (c *Consumer) CloseChan() {
	c.mutex.Lock()
	c.ch.Close()
	c.mutex.Unlock()
}

func (c *Consumer) SetExchangeBinds(eb []*ExchangeBinds) *Consumer {
	c.mutex.Lock()
	if c.state != StateOpened {
		c.exchangeBinds = eb
	}
	c.mutex.Unlock()
	return c
}

func (c *Consumer) SetQueueBinds(qList []string) *Consumer {
	c.mutex.Lock()
	if c.state != StateOpened {
		c.queueList = qList
	}
	c.mutex.Unlock()
	return c
}

// 具体接受的管道信息
func (c *Consumer) SetMsgCallback(cb chan<- Delivery) *Consumer {
	c.mutex.Lock()
	c.callback = cb
	c.mutex.Unlock()
	return c
}

// 设置qos
func (c *Consumer) SetQos(prefetch int) *Consumer {
	c.mutex.Lock()
	c.prefetch = prefetch
	c.mutex.Unlock()
	return c
}

// 该方法需要进行exchange routing的绑定
func (c *Consumer) Open() error {
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
	if c.state == StateOpened {
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
	c.state = StateOpened
	c.stopC = make(chan struct{})
	c.closeC = make(chan *amqp.Error, 1)
	c.ch.NotifyClose(c.closeC)

	opt := DefaultConsumeOption()
	notify := make(chan error, 1)

	c.directConsume(opt, notify)
	for e := range notify {
		if e != nil {
			fmt.Printf("[ERROR] consume has some err %v\n", e)
			continue
		}
		break
	}
	close(notify)

	go c.keepalive()

	return nil
}

// 该模式下直接针对队列名称进行消费，无须绑定exchange
func (c *Consumer) OpenWithDirectConsumeType() error {
	// Open期间不允许对channel做任何操作
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 参数校验
	if c.mq == nil {
		return errors.New("MQ: Bad consumer,can not find any mq")
	}

	// 状态检测
	if c.state == StateOpened {
		return errors.New("MQ: Consumer had been opened")
	}

	// 初始化该consumer的channel，每个消费者独占一个channel
	ch, err := c.mq.channel()
	if err != nil {
		return errors.Errorf("create channel failed %v", err)
	}

	c.ch = ch
	c.state = StateOpened
	c.stopC = make(chan struct{})
	c.closeC = make(chan *amqp.Error, 1)
	c.ch.NotifyClose(c.closeC)

	opt := DefaultConsumeOption()
	notify := make(chan error, 1)

	c.directConsume(opt, notify)
	for e := range notify {
		if e != nil {
			fmt.Printf("[ERROR] consume has some err %v\n", e)
			continue
		}
		break
	}
	close(notify)

	go c.keepalive()

	return nil
}

// 消费者主动关闭
func (c *Consumer) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	select {
	case <-c.stopC:
		// had been closed
	default:
		close(c.stopC)
	}
}

// notifyErr 向上层抛出错误, 如果error为空表示执行完成.由上层负责关闭channel
func (c *Consumer) consume(opt *ConsumeOption, notifyErr chan<- error) {
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
	notifyErr <- nil
}

// 直接进行消费
func (c *Consumer) directConsume(opt *ConsumeOption, notifyErr chan<- error) {
	for i := range c.queueList {
		if len(c.queueList[i]) == 0 {
			notifyErr <- fmt.Errorf("MQ: Queue[%s] is nil, ExchangeBinds[%d], consumer(%s)", c.queueList[i], i, c.name)
			continue
		}
		// 监听消费者绑定的队列，获取通道中的消费信息
		delivery, err := c.ch.Consume(c.queueList[i], "", opt.AutoAck, opt.Exclusive, opt.NoLocal, opt.NoWait, opt.Args)
		if err != nil {
			notifyErr <- fmt.Errorf("MQ: Consumer(%s) consume queue(%s) failed, %v", c.name, c.queueList[i], err)
			continue
		}

		go c.deliver(delivery)
	}

	notifyErr <- nil
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
func (c *Consumer) keepalive() {
	select {
	case <-c.stopC:
		fmt.Printf("[WARN] consumer shutdown normally")
		c.mutex.Lock()
		c.ch.Close()
		c.ch = nil
		c.state = StateClosed
		c.mutex.Unlock()

	case err := <-c.closeC:
		if err == nil {
			fmt.Printf("[ERROR] MQ: Consumer(%s)'s channel was closed, but Error detail is nil\n", c.name)
		} else {
			fmt.Printf("[ERROR] MQ: Consumer(%s)'s channel was closed, code:%d, reason:%s\n", c.name, err.Code, err.Reason)
		}

		// channel被异常关闭了
		c.mutex.Lock()
		c.state = StateReopening
		c.mutex.Unlock()

		var maxRetry = 100
		for i := 0; i < maxRetry; i++ {
			time.Sleep(3 * time.Second)
			// 1. 判断mq是否处于开启状态
			if c.mq.State() != StateOpened {
				fmt.Printf("[WARN] MQ: Consumer(%s) try to recover channel for %d times, but mq's state != StateOpened\n", c.name, i+1)
				continue
			}

			// 2. consumer 打开链接
			e := c.Open()
			if e != nil {
				fmt.Printf("[WARN] MQ: Consumer(%s) recover channel failed for %d times, Err:%v\n", c.name, i+1, e)
				continue
			}

			// 3. 重新打开成功
			fmt.Printf("[INFO] MQ: Consumer(%s) recover channel OK. Total try %d times\n", c.name, i+1)
			return
		}
		// 超过最大重启次数
		fmt.Printf("[ERROR] MQ: Consumer(%s) try to recover channel over maxRetry(%d), so exit\n", c.name, maxRetry)
	}
}

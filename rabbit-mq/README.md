
# rabbit mq


## 简介

rabbit mq是常用的消息中间件，其主要基于 `AMQP` 协议，其核心是 `Queue` `Exchange` `Binding` 为构成AMQP协议的核心。其架构组成如下图：

![img.png](img.png)


- Producer：消息生产者，即投递消息的程序。
- Broker：消息队列服务器实体。
  - Exchange：消息交换机，它指定消息按什么规则，路由到哪个队列。
  - Binding：绑定，它的作用就是把 Exchange 和 Queue 按照路由规则绑定起来。
  - Queue：消息队列载体，每个消息都会被投入到一个或多个队列。
- Consumer：消息消费者，即接受消息的程序。


### Exchange

Exchange接收到信息后，如何将消息转发到对应的Queue中？

根据RoutingKey和当前Exchange所绑定的Binding做匹配。若满足匹配，就向Exchange的Queue中发送消息。Exchange主要有 `Fanout`、`Direct` 和 `Topic` 三种类型。由于本次业务中主要使用了 `Direct` 类型，所以这里重点就Direct进行demo的编写和演示。





## 参考

1. [rabbitmq exchange讲解](https://zhuanlan.zhihu.com/p/37198933)
package subscriber

import (
	"fmt"
	"github.com/streadway/amqp"
	"sync"
)

type Subscriber struct {
	host       string
	connection *amqp.Connection
	handler    Handler
}

type Handler interface {
	Handle([]byte)
}

// NewSubscriber 创建一个Subscriber并且返回其指针
func NewSubscriber(host string) (*Subscriber, error) {
	conn, _ := amqp.Dial(host)
	sub := new(Subscriber)
	sub.host = host
	sub.connection = conn
	return sub, nil
}

// Subscribe 将队列与指定交换机绑定并开始监听，传入参数为队列名称、
func (p *Subscriber) Subscribe(exchangeName string, handler Handler) error {
	ch, err := p.connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName,
		amqp.ExchangeFanout,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return err
	}

	// 使用临时队列，不指定名称而是自动生成
	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	err = ch.QueueBind(
		q.Name,       // queue name
		"",           // routing key
		exchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	// 处理队列中消息的协程
	go func() {
		for d := range msgs {
			//fmt.Println("Get msg now")
			// 可根据d.contentType选择不同的处理函数
			handler.Handle(d.Body)
		}
		fmt.Println("forever协程运行到要结束的地方咯")
		forever <- true
	}()

	<-forever
	println("stop subscribe!!!")
	return nil
}

// SubscribeWithChan 将队列与指定交换机绑定并开始监听，传入参数为队列名称
func (p *Subscriber) SubscribeWithSync(exchangeName string, handler Handler, wg *sync.WaitGroup) error {
	ch, err := p.connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName,
		amqp.ExchangeFanout,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return err
	}

	// 使用临时队列，不指定名称而是自动生成
	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	err = ch.QueueBind(
		q.Name,       // queue name
		"",           // routing key
		exchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	// 处理队列中消息的协程
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			//fmt.Println("Get msg now")
			// 可根据d.contentType选择不同的处理函数
			handler.Handle(d.Body)
		}
		fmt.Println("协程运行到要结束的地方咯")
		forever <- true
	}()
	<-forever
	wg.Done()
	return nil
}

// CloseConnection 关闭到消息队列的连接
func (p *Subscriber) CloseConnection() error {
	if !p.connection.IsClosed() {
		err := p.connection.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

package publisher

import (
	"github.com/streadway/amqp"
)

type Publisher struct {
	connection *amqp.Connection
	host       string
}

// NewPublisher 创建一个Publisher并且返回其指针
func NewPublisher(host string) (*Publisher, error) {
	conn, _ := amqp.Dial(host)
	pub := new(Publisher)
	pub.host = host

	pub.connection = conn
	return pub, nil
}

// Publish 向指定的交换机广播一条信息并立即返回，广播类型为FANOUT
func (p *Publisher) Publish(queueName string, body []byte, contentType string) error {
	ch, err := p.connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		queueName,
		amqp.ExchangeFanout,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return err
	}

	err = ch.Publish(
		queueName,
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: contentType,
			Body:        body,
		})
	if err != nil {
		return err
	}
	return nil
}

// CloseConnection 关闭到消息队列的连接
func (p *Publisher) CloseConnection() error {
	if !p.connection.IsClosed() {
		err := p.connection.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

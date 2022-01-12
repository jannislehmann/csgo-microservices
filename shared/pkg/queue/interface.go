package queue

import "github.com/streadway/amqp"

type QueueUseCase interface {
	Connect(amqpURI string) error
	CreateQueue(topic string) error
	Publish(body []byte, topic string) (<-chan bool, error)
	Consume(topic string) (<-chan amqp.Delivery, error)
}

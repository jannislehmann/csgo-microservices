package queue

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type QueueService struct {
	Channel    *amqp.Channel
	Connection *amqp.Connection
}

func NewService() *QueueService {
	return &QueueService{}
}

func (s *QueueService) Connect(amqpURI string) error {
	connection, err := amqp.Dial(amqpURI)
	if err != nil {
		return fmt.Errorf("dial: %s", err)
	}
	s.Connection = connection

	return nil
}

func (s *QueueService) CreateQueue(topic string) error {
	channel, err := s.Connection.Channel()
	if err != nil {
		return fmt.Errorf("channel: %s", err)
	}

	if _, err := channel.QueueDeclare(
		topic,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("queue: %s", err)
	}

	log.Printf("enabling publishing confirms.")
	if err := channel.Confirm(false); err != nil {
		return fmt.Errorf("channel could not be put into confirm mode: %s", err)
	}

	s.Channel = channel

	return nil
}

func (s *QueueService) Publish(body []byte, topic string) (<-chan bool, error) {
	if err := s.Channel.Publish(
		"",
		topic,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	); err != nil {
		return nil, fmt.Errorf("publish: %s", err)
	}

	// Send acknowledge to channel.
	confirms := s.Channel.NotifyPublish(make(chan amqp.Confirmation, 1))
	resultChannel := make(chan bool, 1)
	defer func(confirmChannel chan amqp.Confirmation, ackChannel chan bool) {
		confirmed := <-confirmChannel
		ackChannel <- confirmed.Ack
	}(confirms, resultChannel)

	return resultChannel, nil
}

func (s *QueueService) Consume(topic string) (<-chan amqp.Delivery, error) {
	return s.Channel.Consume(
		topic,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
}

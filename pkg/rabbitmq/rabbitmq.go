package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

func OpenChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	return ch
}

func Consume(ch *amqp.Channel, out chan amqp.Delivery, queue, consumerName string) error {
	msgs, err := ch.Consume(
		queue,
		consumerName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	for msg := range msgs {
		out <- msg
	}
	return nil
}

func Publish(ch *amqp.Channel, msg amqp.Publishing, exName string) error {
	err := ch.PublishWithContext(
		context.Background(),
		exName,
		"",
		false,
		false,
		msg,
	)
	if err != nil {
		return err
	}
	return nil
}

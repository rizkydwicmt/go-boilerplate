package rabbitmq

import (
	helper "go-boilerplate/utilities"

	"github.com/rabbitmq/amqp091-go"
)

func RabbitMQ() (*amqp091.Connection, *amqp091.Channel) {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		helper.Log("error")("failed to connect to amqp", err)
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		helper.Log("error")("failed to get channel", err)
		panic(err)
	}

	return conn, ch
}

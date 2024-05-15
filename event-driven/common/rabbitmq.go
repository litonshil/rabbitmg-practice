package common

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

const OrderCreateEvent = "order.created"

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
func ConnectAmqp(user, pass, host, port string) (*amqp.Channel, func() error) {
	connectString := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pass, host, port)
	conn, err := amqp.Dial(connectString)
	FailOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")

	err = ch.ExchangeDeclare(
		OrderCreateEvent, // name
		"fanout",         // type
		true,             // durable
		false,            // auto-deleted
		false,            // internal
		false,            // no-wait
		nil,              // arguments
	)
	FailOnError(err, "Failed to declare an exchange")
	return ch, conn.Close
}

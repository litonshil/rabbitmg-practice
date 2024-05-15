package main

import (
	"context"
	"encoding/json"
	common "github.com/litonshil009/rabbitmq-event-driven/common"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

var (
	amqpUser = "guest"
	amqpPass = "guest"
	amqpHost = "localhost"
	amqpPort = "5672"
)

func main() {
	ch, close := common.ConnectAmqp(amqpUser, amqpPass, amqpHost, amqpPort)
	defer func() {
		close()
		ch.Close()
	}()

	q, err := ch.QueueDeclare(common.OrderCreateEvent, true, false, false, false, nil)
	common.FailOnError(err, "error at queue declare")

	orderReq, err := json.Marshal(common.Order{
		ID: "order-1",
		Items: []common.Item{
			{
				ID:       "item-1",
				Quantity: 1,
			},
		},
	})

	common.FailOnError(err, "failed to marshal order payload")
	err = ch.PublishWithContext(context.Background(), "", q.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        orderReq,
	})
	common.FailOnError(err, "Failed to publish a message")
	log.Println("Order Published")
}

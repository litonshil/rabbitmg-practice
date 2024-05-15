package main

import (
	"encoding/json"
	"github.com/litonshil009/rabbitmq-event-driven/common"
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

	listen(ch)
}

func listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(common.OrderCreateEvent, true, false, false, false, nil)
	common.FailOnError(err, "error at decalring queue")

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	common.FailOnError(err, "failed to consume")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			order := &common.Order{}
			if err := json.Unmarshal(d.Body, order); err != nil {
				d.Nack(false, false)
				log.Printf("failed to unmarshal order: %v", err)
				continue
			}

			paymentLink, err := createPaymentLink()
			if err != nil {
				common.FailOnError(err, "failed to create payment")
				// handle retry.....
				continue
			}
			log.Printf("payment link genreated: %s", paymentLink)
		}

	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}

func createPaymentLink() (string, error) {
	return "dummy-payment-link.com", nil
}

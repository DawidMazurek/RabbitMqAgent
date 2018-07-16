package main

import "github.com/streadway/amqp"

func declareExchange(ch *amqp.Channel, obj Message) {
	err := ch.ExchangeDeclare(
		obj.Deliver_options.Exchange_name,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to declare exchange")
}

func publishMessage(ch *amqp.Channel, obj Message) {
	err := ch.Publish(
		obj.Deliver_options.Exchange_name,
		obj.Deliver_options.Routing_key,
		false,
		false,
		amqp.Publishing {
			ContentType: "text/plain",
			Body:        []byte(obj.Payload),
		})

	failOnError(err, "Failed to publish a message")
}
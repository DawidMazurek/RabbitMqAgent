package main

import (
	"net"
	"encoding/json"
	"fmt"
	"log"
	"github.com/streadway/amqp"
	"os"
)

func deliverMessage(c net.Conn) {
	for {
		buf := make([]byte, 512)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}

		data := buf[0:nr]
		println("Received: ", string(data))

		enqueueToRabbitMQ(data)
	}
}

func dumpDetails(obj Message) {
	fmt.Println("Sending to rabbitmq")
	fmt.Printf("vhost: %s\n", obj.Deliver_options.Vhost)
	fmt.Printf("echange: %s\n", obj.Deliver_options.Exchange_name)
	fmt.Printf("routing key: %s\n", obj.Deliver_options.Routing_key)
	fmt.Printf("Writing response to socket: %s\n", obj.Response_socket)
}

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

func enqueueToRabbitMQ(data []byte) {
	var obj Message
	err := json.Unmarshal(data, &obj)

	dumpDetails(obj)

	rabbitDSN := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASS"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
		obj.Deliver_options.Vhost,
	)

	conn, err := amqp.Dial(rabbitDSN)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()


	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	declareExchange(ch, obj)

	err = ch.Publish(
		obj.Deliver_options.Exchange_name,
		obj.Deliver_options.Routing_key,
		false,
		false,
		amqp.Publishing {
			ContentType: "text/plain",
			Body:        []byte(obj.Payload),
		})

	failOnError(err, "Failed to publish a message")

	//////////////////////////////////////////////////

	var responseSocket string = obj.Response_socket

	c, err := net.Dial("unix", responseSocket)

	if err != nil {
		panic(err)
	}
	_, err = c.Write([]byte("response"))

	if err != nil {
		log.Println(err)
	}
	fmt.Println("Sent.")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
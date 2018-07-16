package main

import (
	"net"
	"encoding/json"
	"fmt"
	"log"
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

func enqueueToRabbitMQ(data []byte) {
	var obj Message
	err := json.Unmarshal(data, &obj)

	dumpDetails(obj)

	conn := getConnection(obj.Deliver_options.Vhost);

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	declareExchange(ch, obj)
	publishMessage(ch, obj)
	sendResponseToSocket(obj)
}

func sendResponseToSocket(obj Message) {
	var responseSocket = obj.Response_socket

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

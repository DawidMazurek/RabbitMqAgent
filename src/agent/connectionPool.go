package main

import (
	"github.com/streadway/amqp"
	"fmt"
	"os"
)

var connectionPool map[string] *amqp.Connection

func getConnection(vhost string) *amqp.Connection {
	var conn *amqp.Connection
	var connExists bool
	conn, connExists = connectionPool[vhost]
	if connExists {
		return conn
	}

	rabbitDSN := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASS"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
		vhost,
	)

	conn, err := amqp.Dial(rabbitDSN)
	failOnError(err, "Failed to connect to RabbitMQ")

	if connectionPool == nil {
		connectionPool = map[string] *amqp.Connection{}
	}

	connectionPool[vhost] = conn

	return conn
}

package main

import (
	"github.com/streadway/amqp"
	"fmt"
)

var connectionPool map[string] *amqp.Connection

func getConnection(vhost string) *amqp.Connection {
	var conn *amqp.Connection
	var connExists bool
	conn, connExists = connectionPool[vhost]
	if connExists {
		return conn
	}

	config := getConnectionConfig()

	rabbitDSN := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		config.Get("user"),
		config.Get("pass"),
		config.Get("host"),
		config.Get("port"),
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

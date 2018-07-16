package main

type DeliverOptions struct {
	Vhost string
	Exchange_name string
	Routing_key string
}

type Message struct {
	Payload string
	Deliver_options DeliverOptions
	Response_socket string
}

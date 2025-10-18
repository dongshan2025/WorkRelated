package main

import (
	"flag"

	"github.com/streadway/amqp"
)

var url = flag.String("url", "amqp://guest:guest@localhost:5672/", "amqp url for the publisher and subscriber")

// exchange binds the publishers to the subscribers
const exchange = "pubsub"

// message is the application type for a message
// this can contain identity, or a reference to the receiver chan for further demuxing
type message []byte

// session composes an amqp.Connection with an amqp.Channel
type session struct {
	*amqp.Connection
	*amqp.Channel
}

// close tears the connection down, taking the channel with it
func (s session) Close() error {
	if s.Connection == nil {
		return nil
	}
	return s.Connection.Close()
}

// redial continually connects to the URL, exiting the program when no longer possible

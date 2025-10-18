package main

import (
	"log"

	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("hello", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	body := "Hello World"
	err = ch.Publish("", q.Name, false, false, amqp.Publishing{ContentType: "text/plain", Body: []byte(body)})
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("[x] Sent %s\n", body)
}

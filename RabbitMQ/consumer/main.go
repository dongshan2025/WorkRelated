package main

import (
	"log"

	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	go func() {
		for d := range msgs {
			log.Printf("Receive a message: %s", d.Body)
		}
	}()

	log.Printf("[*]Waiting for messages. To exit press CTRL+C")
	select {}
}

package main

import (
	"log"
	"strconv"

	"github.com/streadway/amqp"
)

func fib(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return fib(n-1) + fib(n-2)
	}
}

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

	q, err := ch.QueueDeclare("rpc_queue", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = ch.Qos(1, 0, false)
	if err != nil {
		log.Fatal(err)
		return
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	go func() {
		for d := range msgs {
			n, err := strconv.Atoi(string(d.Body))
			if err != nil {
				log.Fatal(err)
				return
			}

			log.Printf("[.] fib(%d)", n)

			response := fib(n)

			err = ch.Publish("", d.ReplyTo, false, false, amqp.Publishing{ContentType: "text/plain", CorrelationId: d.CorrelationId, Body: []byte(strconv.Itoa(response))})
			if err != nil {
				log.Fatal(err)
				return
			}

			d.Ack(false)
		}
	}()

	log.Printf("[*] Awaiting RPC requests")
	select {}
}

package main

import (
	"log"
	"math/rand"
	"strconv"

	"github.com/streadway/amqp"
)

func main() {
	res, err := fibonacciRPC(6)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("[.] Got %d", res)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func randomString(n int) string {
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func fibonacciRPC(n int) (int, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return 0, err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		return 0, err
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return 0, err
	}

	corrid := randomString(32)

	err = ch.Publish("", "rpc_queue", false, false, amqp.Publishing{ContentType: "text/plain", CorrelationId: corrid, ReplyTo: q.Name, Body: []byte(strconv.Itoa(n))})
	if err != nil {
		return 0, err
	}

	for d := range msgs {
		if corrid == d.CorrelationId {
			res, err := strconv.Atoi(string(d.Body))
			if err != nil {
				return 0, err
			}

			return res, nil
		}
	}

	return 0, nil
}

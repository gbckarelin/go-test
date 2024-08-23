package main

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatal("%s: %s", msg, err)
	}
}

func StartConsumer() {
	var conn *amqp.Connection
	var err error

	for ret := 0; ret < 10; ret++ {
		conn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err == nil {
			break
		}
		log.Printf("Failed to connect to RabbitMQ. Retrying in 10 seconds...")
		time.Sleep(10 * time.Second)
	}

	ch, err := conn.Channel()
	failOnError(err, "Failed to open channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"qqq", // name
		false, // passive
		false, // durable
		false, //autodelete
		false, //nowait
		nil,   //args
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			taskID := string(d.Body)
			store.doTask(taskID)
			log.Printf("Recived a message: %s", taskID)
		}
	}()

	log.Printf("[*] Waiting for messages. To exit press Ctrl+C")
	<-forever
}

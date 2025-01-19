package rabbitmq

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQService interface {
	SendCSVToQueueue(idx int, data []string) error
	ReceiveFromQueue() (<-chan amqp.Delivery, error)
}

func ConnectToRabbitMQ() (*RabbitMQ, error) {
	RABBIT_MQ_CONNECTION_STRING := os.Getenv("RABBIT_MQ_CONNECTION_STRING")
	QUEUE_NAME := os.Getenv("QUEUE_NAME")
	conn, err := amqp.Dial(RABBIT_MQ_CONNECTION_STRING)
	if err != nil {
		log.Println("Error in connecting to rabbitmq", err)
		return nil, err
	}

	//defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Println("Error in creating channel", err)
		return nil, err
	}
	//defer ch.Close()

	q, err := ch.QueueDeclare(
		QUEUE_NAME,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Println("Error in declaring queue", err)
		return nil, err
	}

	return &RabbitMQ{Queue: &q, Channel: ch}, nil

}

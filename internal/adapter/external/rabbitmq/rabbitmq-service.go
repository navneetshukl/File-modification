package rabbitmq

import (
	"file-modification/internal/adapter/external/csv"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQService interface {
	SendCSVToQueueue(fileName string) error
}

type QueueStruct struct {
	Queue   *amqp.Queue
	Channel *amqp.Channel
}

func ConnectToRabbitMQ() (*QueueStruct, error) {
	RABBIT_MQ_CONNECTION_STRING := os.Getenv("RABBIT_MQ_CONNECTION_STRING")
	QUEUE_NAME := os.Getenv("QUEUE_NAME")
	conn, err := amqp.Dial(RABBIT_MQ_CONNECTION_STRING)
	if err != nil {
		log.Println("Error in connecting to rabbitmq", err)
		return nil, err
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Println("Error in creating channel", err)
		return nil, err
	}
	defer ch.Close()

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

	return &QueueStruct{Queue: &q, Channel: ch}, nil

}

func NewRabbitMQService(csv csv.CSVService) (*RabbitMQ, error) {
	rabbitStruct, err := ConnectToRabbitMQ()
	if err != nil {
		log.Println("Error in connecting to rabbitmq", err)
		return nil, err
	}
	queueStruct := &QueueStruct{Queue: rabbitStruct.Queue, Channel: rabbitStruct.Channel}
	return &RabbitMQ{RabbitMQStruct: *queueStruct, CSRService: csv}, nil

}

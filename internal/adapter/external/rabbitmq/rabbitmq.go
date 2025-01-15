package rabbitmq

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Queue   *amqp.Queue
	Channel *amqp.Channel
}

func (r *RabbitMQ) SendCSVToQueueue(data []string) error {

	dataByte, err := json.Marshal(data)
	if err != nil {
		log.Println("error in converting to byte ", err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = r.Channel.PublishWithContext(ctx,
		"",
		r.Queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        dataByte,
		})
	if err != nil {
		log.Println("failed to publish a message:", err)
		return err
	}

	return nil

}

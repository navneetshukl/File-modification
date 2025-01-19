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

type queueData struct {
	OpType   int      `json:"op_type"`
	Sequence int      `json:"sequence"`
	Data     []string `json:"data"`
}

func (r *RabbitMQ) SendCSVToQueueue(seq int, csvdata []string) error {

	qdata := &queueData{
		OpType:   1,
		Sequence: seq,
		Data:     csvdata,
	}

	dataByte, err := json.Marshal(qdata)
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

func (r *RabbitMQ) ReceiveFromQueue() (<-chan amqp.Delivery, error) {
	msgs, err := r.Channel.Consume(
		r.Queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("error in consuming from queue ", err)
		return nil, err
	}
	return msgs, nil

}

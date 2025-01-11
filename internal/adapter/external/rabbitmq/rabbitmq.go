package rabbitmq

import (
	"file-modification/internal/adapter/external/csv"
	"log"
)

const LineChunk int = 1000

type RabbitMQ struct {
	RabbitMQStruct QueueStruct
	CSRService     csv.CSVService
}

func (r *RabbitMQ) SendCSVToQueueue(fileName string) error {

	_, err := r.CSRService.ReadCSV(fileName)
	if err != nil {
		log.Println("Error in getting the csv data ", err)
		return err
	}

	return nil

}

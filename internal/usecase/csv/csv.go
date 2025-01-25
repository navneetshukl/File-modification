package csv

import (
	"context"
	"encoding/json"
	"file-modification/internal/adapter/external/csv"
	"file-modification/internal/adapter/external/rabbitmq"
	csvCore "file-modification/internal/core/csv"
	"fmt"
	"log"
	"sync"
)

type CsvServiceImpl struct {
	CsvReaderSvc csv.CSVService
	RabbitSvc    rabbitmq.RabbitMQService
	wg           *sync.WaitGroup
	size         int
}

func NewCsvUseCaseImpl(csv csv.CSVService, rabb rabbitmq.RabbitMQService) *CsvServiceImpl {
	return &CsvServiceImpl{
		CsvReaderSvc: csv,
		RabbitSvc:    rabb,
		wg:           &sync.WaitGroup{},
	}
}

func (p *CsvServiceImpl) ReadCSV(ctx context.Context, fileName string) error {

	str, err := p.CsvReaderSvc.ReadCSV(fileName)
	if err != nil {
		log.Printf("Error in reading the pdf %s.Error is %v\n", fileName, err)
		return err

	}
	p.size = len(str)
	fmt.Println("size is ", p.size)
	for idx, val := range str {
		err := p.RabbitSvc.SendCSVToQueueue(idx+1, val)
		if err != nil {
			log.Printf("Error in sending the pdf %s.Error is %v\n", fileName, err)
			return err
		}
	}

	p.processCSV()

	return nil
}

func (p *CsvServiceImpl) processCSV() error {
	msg, err := p.RabbitSvc.ReceiveFromQueue()
	if err != nil {
		log.Println("error in reading from the queue ", err)
		return err

	}

	for d := range msg {
		data := csvCore.CSVData{}
		err := json.Unmarshal(d.Body, &data)
		if err != nil {
			log.Println("error in unmarshalling the data ", err)
		}
		fmt.Println("Received a message: ", data)
		if data.Sequence == p.size {
			break
		}

	}

	return nil
}

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

	"github.com/rabbitmq/amqp091-go"
)

type CsvServiceImpl struct {
	CsvReaderSvc csv.CSVService
	RabbitSvc    rabbitmq.RabbitMQService
	wg           *sync.WaitGroup
	size         int
	numWorkers   int
	taksk        int
}

func NewCsvUseCaseImpl(csv csv.CSVService, rabb rabbitmq.RabbitMQService) *CsvServiceImpl {
	return &CsvServiceImpl{
		CsvReaderSvc: csv,
		RabbitSvc:    rabb,
		wg:           &sync.WaitGroup{},
		numWorkers:   5,
	}
}

func (p *CsvServiceImpl) ReadCSV(ctx context.Context, fileName string) error {

	str, err := p.CsvReaderSvc.ReadCSV(fileName)
	if err != nil {
		log.Printf("Error in reading the pdf %s.Error is %v\n", fileName, err)
		return err

	}
	p.size = len(str)
	p.taksk = p.size
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

	tasks := make(chan amqp091.Delivery, p.numWorkers)
	for w := 1; w <= p.numWorkers; w++ {
		p.wg.Add(1)
		go p.processEachRow(context.Background(), tasks, p.wg)
	}

	for d := range msg {
		tasks <- d
	}
	p.wg.Wait()
	close(tasks)

	return nil
}

func (p *CsvServiceImpl) processEachRow(ctx context.Context, tasks <-chan amqp091.Delivery, wg *sync.WaitGroup) {
	defer wg.Done()

	for d := range tasks {
		data := csvCore.CSVData{}
		err := json.Unmarshal(d.Body, &data)
		if err != nil {
			log.Println("error in unmarshalling the data ", err)
		}
		if data.OperationType == 1 {
			fmt.Println("Received a message: ", data)

			// do some task
		}
	}

}

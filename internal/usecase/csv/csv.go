package csv

import (
	"context"
	"encoding/json"
	"file-modification/internal/adapter/external/csv"
	"file-modification/internal/adapter/external/rabbitmq"
	csvCore "file-modification/internal/core/csv"
	"fmt"
	"log"
	"strings"
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
	pdf          [][]string
}

func NewCsvUseCaseImpl(csv csv.CSVService, rabb rabbitmq.RabbitMQService) *CsvServiceImpl {
	return &CsvServiceImpl{
		CsvReaderSvc: csv,
		RabbitSvc:    rabb,
		wg:           &sync.WaitGroup{},
		numWorkers:   5,
		pdf:          [][]string{},
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
		err := p.RabbitSvc.SendCSVToQueueue(1, idx+1, val)
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
		go p.processEachRow(tasks)
	}

	for i := 0; i < p.size; i++ {
		tasks <- <-msg
	}

	close(tasks)
	p.wg.Wait()
	fmt.Println("Total pdf ", len(p.pdf))
	return nil
}

func (p *CsvServiceImpl) processEachRow(tasks <-chan amqp091.Delivery) {
	defer p.wg.Done()

	for d := range tasks {
		data := csvCore.CSVData{}
		err := json.Unmarshal(d.Body, &data)
		if err != nil {
			log.Println("error in unmarshalling the data ", err)
		}
		if data.OperationType == 1 {
			str := []string{}
			fmt.Println("Inside operation type 1")
			for _, val := range data.Data {
				s := strings.ToUpper(val)
				str = append(str, s)
			}
			// err := p.RabbitSvc.SendCSVToQueueue(2, data.Sequence, str)
			// if err != nil {
			// 	log.Println("error in sending the data to the queue ", err)
			// 	return
			// }

			p.pdf = append(p.pdf, str)
		}
	}
}

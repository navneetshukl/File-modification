package csv

import (
	"context"
	"file-modification/internal/adapter/external/csv"
	"file-modification/internal/adapter/external/rabbitmq"
	"log"
	"sync"
)

type CsvServiceImpl struct {
	CsvReaderSvc csv.CSVService
	RabbitSvc    rabbitmq.RabbitMQService
	wg           *sync.WaitGroup
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
	for idx, val := range str {
		err := p.RabbitSvc.SendCSVToQueueue(idx+1, val)
		if err != nil {
			log.Printf("Error in sending the pdf %s.Error is %v\n", fileName, err)
			return err
		}
	}

	// msg, err := p.RabbitSvc.ReceiveFromQueue()
	// if err != nil {
	// 	log.Printf("Error in reading the pdf %s.Error is %v\n", fileName, err)
	// 	return err

	// }
	// ch := make(chan int)

	// go func() {
	// 	defer close(ch)
	// 	for d := range msg {
	// 		log.Println("Received message from queue ", string(d.Body))

	// 	}

	// }()

	// for val:=range ch {
	// 	log.Println("val is ",val)
	// }

	return nil
}

func (p *CsvServiceImpl) processCSV() error {
	msg, err := p.RabbitSvc.ReceiveFromQueue()
	if err != nil {
		log.Println("error in reading from the queue ", err)
		return err

	}

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for d := range msg {
			log.Println("Received message from queue ", string(d.Body))

		}

	}()
	p.wg.Wait()

	return nil
}

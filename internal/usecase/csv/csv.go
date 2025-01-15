package csv

import (
	"file-modification/internal/adapter/external/csv"
	"file-modification/internal/adapter/external/rabbitmq"
	"log"
	"sync"
)

type CsvServiceImpl struct {
	CsvReaderSvc csv.CSVService
	RabbitSvc    rabbitmq.RabbitMQService
	wg           sync.WaitGroup
}

func NewCsvUseCaseImpl(csv csv.CSVService, rabb rabbitmq.RabbitMQService) *CsvServiceImpl {
	return &CsvServiceImpl{
		CsvReaderSvc: csv,
		RabbitSvc:    rabb,
		wg:           sync.WaitGroup{},
	}
}

func (p *CsvServiceImpl) ReadCSV(fileName string) error {

	str, err := p.CsvReaderSvc.ReadCSV(fileName)
	if err != nil {
		log.Printf("Error in reading the pdf %s.Error is %v\n", fileName, err)
		return err

	}

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for _, val := range str {
			err := p.RabbitSvc.SendCSVToQueueue(val)
			if err != nil {
				log.Printf("Error in sending to RabbitMQ %v,%v :", val, err)

			}

		}

	}()
	p.wg.Wait()
	return nil
}

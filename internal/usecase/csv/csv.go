package csv

import (
	"file-modification/internal/adapter/external/csv"
	"file-modification/internal/adapter/external/rabbitmq"
	"log"
)

type CsvServiceImpl struct {
	CsvReaderSvc csv.CSVService
	RabbitSvc    rabbitmq.RabbitMQService
}

func NewCsvUseCaseImpl(csv csv.CSVService, rabb rabbitmq.RabbitMQService) *CsvServiceImpl {
	return &CsvServiceImpl{
		CsvReaderSvc: csv,
		RabbitSvc:    rabb,
	}
}

func (p *CsvServiceImpl) ReadCSV(fileName string) error {

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

	msg, err := p.RabbitSvc.ReceiveFromQueue()
	if err != nil {
		log.Printf("Error in reading the pdf %s.Error is %v\n", fileName, err)
		return err

	}
	ch := make(chan bool)

	go func() {
		for d := range msg {
			log.Println("Received message from queue ", string(d.Body))
		}
		ch<-true
	}()

	<-ch

	return nil
}

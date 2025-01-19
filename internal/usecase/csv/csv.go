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

	msg,err:=p.RabbitSvc.ReceiveFromQueue()
	if err != nil {
		log.Printf("Error in reading the pdf %s.Error is %v\n", fileName, err)
		return err

	}

	for d := range msg {
		log.Printf("Received a message: %s", d.Body)
	}
	return nil
}

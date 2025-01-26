package csv

import (
	"context"
	"encoding/json"
	"file-modification/internal/adapter/external/csv"
	"file-modification/internal/adapter/external/pdf"
	"file-modification/internal/adapter/external/rabbitmq"
	csvCore "file-modification/internal/core/csv"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"

	"github.com/rabbitmq/amqp091-go"
)

type CsvServiceImpl struct {
	CsvReaderSvc csv.CSVService
	RabbitSvc    rabbitmq.RabbitMQService
	PdfSvc       pdf.PDFService
	wg           *sync.WaitGroup
	size         int
	numWorkers   int
	taksk        int
	pdf          []csvCore.CSVData
}

func NewCsvUseCaseImpl(csv csv.CSVService, rabb rabbitmq.RabbitMQService, pdf pdf.PDFService) *CsvServiceImpl {
	return &CsvServiceImpl{
		CsvReaderSvc: csv,
		RabbitSvc:    rabb,
		wg:           &sync.WaitGroup{},
		numWorkers:   5,
		pdf:          make([]csvCore.CSVData, 0),
		PdfSvc:       pdf,
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

	err = p.convertToPDF()
	if err != nil {
		log.Println("error in converting to pdf ", err)
		return err
	}

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

			pdfData := csvCore.CSVData{
				OperationType: 2,
				Sequence:      data.Sequence,
				Data:          str,
			}

			p.pdf = append(p.pdf, pdfData)
		}
	}
}

func (p *CsvServiceImpl) convertToPDF() error {
	sort.Slice(p.pdf, func(i, j int) bool {
		return p.pdf[i].Sequence < p.pdf[j].Sequence
	})

	fmt.Println("pdf data is ", p.pdf[0].Data)

	pdfData := [][]string{}

	for _, val := range p.pdf {
		pdfData = append(pdfData, val.Data)
	}

	err := p.PdfSvc.ConvertToPDF("data.pdf", pdfData)
	if err != nil {
		log.Println("error in converting to pdf ", err)
		return err
	}
	return nil

}

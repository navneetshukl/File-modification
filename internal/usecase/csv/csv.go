package csv

import (
	"file-modification/internal/adapter/external/csv"
	"fmt"
	"log"
)

type CsvServiceImpl struct {
	CsvReader csv.CSVService
}

func NewCsvUseCaseImpl(csv csv.CSVService) *CsvServiceImpl {
	return &CsvServiceImpl{
		CsvReader: csv,
	}
}

func (p *CsvServiceImpl) ReadCSV(fileName string) error {
	fileName = "uploads/resume.pdf"

	str, err := p.CsvReader.ReadCSV(fileName)
	if err != nil {

		log.Printf("Error in reading the pdf %s.Error is %v\n", fileName, err)
		return err

	}

	fmt.Println(str)
	return nil
}

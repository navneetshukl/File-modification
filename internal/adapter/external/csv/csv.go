package csv

import (
	"log"
	"os"
)

type CSVService interface {
	ReadCSV(fileName string) ([]*string, error)
}

type CsvServiceImpl struct{}

func NewCSVService() *CsvServiceImpl {
	return &CsvServiceImpl{}
}

func (p *CsvServiceImpl) ReadCSV(fileName string) ([]*string, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Println("Error in getting the current workig directory ", err)
		return nil, err
	}
	fileName = dir + "/uploads/" + fileName
	return nil, nil
}

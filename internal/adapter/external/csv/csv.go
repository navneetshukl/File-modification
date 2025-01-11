package csv

import (
	"encoding/csv"
	"log"
	"os"
)

type CSVService interface {
	ReadCSV(fileName string) ([][]string, error)
}

type CsvServiceImpl struct{}

func NewCSVService() *CsvServiceImpl {
	return &CsvServiceImpl{}
}

func (p *CsvServiceImpl) ReadCSV(fileName string) ([][]string, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Println("Error in getting the current workig directory ", err)
		return nil, err
	}
	fileName = dir + "/uploads/" + fileName

	file, err := os.Open(fileName)
	if err != nil {
		log.Println("Error in opening the file ", err)
		return nil, err
	}

	defer file.Close()

	reader := csv.NewReader(file)

	csvData := [][]string{}

	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Println("Error in reading the csv file ", err)
			return nil, err
		}
		csvData = append(csvData, record)
	}
	return csvData, nil
}

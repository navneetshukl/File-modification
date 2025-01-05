package csv

type CSVService interface {
	ReadCSV(fileName string) (string, error)
}

type CsvServiceImpl struct{}

func NewCSVService() *CsvServiceImpl {
	return &CsvServiceImpl{}
}

func(p *CsvServiceImpl) ReadCSV(fileName string) (string, error) {
	return "",nil
}



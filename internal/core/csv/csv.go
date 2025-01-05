package csv

type CSVUseCase interface {
	ReadCSV(fileName string) error
}

package csv

import "context"

type CSVUseCase interface {
	ReadCSV(ctx context.Context, fileName string) error
}

type CSVData struct {
	OperationType int      `json:"op_type"`
	Sequence      int      `json:"sequence"`
	Data          []string `json:"data"`
}

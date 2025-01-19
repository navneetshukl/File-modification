package csv

import "context"

type CSVUseCase interface {
	ReadCSV(ctx context.Context, fileName string) error
}

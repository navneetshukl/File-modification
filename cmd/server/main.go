package main

import (
	"file-modification/internal/adapter/external/csv"
	s3Service "file-modification/internal/adapter/external/s3"
	routes "file-modification/internal/interface/api"
	"file-modification/internal/interface/api/handler"
	csvImpl "file-modification/internal/usecase/csv"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		return
	}
	_, err = s3Service.NewS3ServiceImpl()
	if err != nil {
		log.Println("Error in creating s3 client")
		return
	}

	csvService := csv.NewCSVService()
	csvUseCase := csvImpl.NewCsvUseCaseImpl(csvService)

	csvUseCase.ReadCSV("data.csv")

	h := handler.NewHandler()
	app := routes.SetupRoutes(h)
	app.Listen(":8080")

}

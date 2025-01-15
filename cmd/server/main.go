package main

import (
	"file-modification/internal/adapter/external/csv"
	"file-modification/internal/adapter/external/rabbitmq"
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
	rabbitService, err := rabbitmq.ConnectToRabbitMQ()
	if err != nil {
		log.Println("Error in connecting to rabbitmq")
		return
	}
	csvUseCase := csvImpl.NewCsvUseCaseImpl(csvService, rabbitService)

	csvUseCase.ReadCSV("data.csv")

	h := handler.NewHandler()
	app := routes.SetupRoutes(h)
	app.Listen(":8080")

}

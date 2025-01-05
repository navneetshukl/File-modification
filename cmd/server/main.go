package main

import (
	"file-modification/internal/adapter/external/pdf"
	s3Service "file-modification/internal/adapter/external/s3"
	routes "file-modification/internal/interface/api"
	pdfImpl"file-modification/internal/usecase/pdf"
	"file-modification/internal/interface/api/handler"
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

	pdfService:=pdf.NewPDFService()
	pdfUsecase:=pdfImpl.NewPdfServiceImpl(pdfService)

	pdfUsecase.ReadPDF("scholarship.pdf")


	h := handler.NewHandler()
	app := routes.SetupRoutes(h)
	app.Listen(":8080")

}

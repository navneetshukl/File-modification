package main

import (
	s3Service "file-modification/internal/adapter/external/s3"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	err:=godotenv.Load()
	if err!=nil{
		log.Println("Error loading .env file")
		return
	}
	s3Service.NewS3ServiceImpl()
}

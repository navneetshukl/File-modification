package main

import (
	s3Service "file-modification/internal/adapter/external/s3"
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

}

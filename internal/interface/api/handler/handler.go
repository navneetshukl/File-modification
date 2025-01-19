package handler

import (
	"file-modification/internal/core/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	csvUseCase csv.CSVUseCase
}

func NewHandler(csv csv.CSVUseCase) *Handler {
	return &Handler{
		csvUseCase: csv,
	}
}

func (h *Handler) UploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to upload file",
		})
	}

	// Check if the file size exceeds 10 MB (10 * 1024 * 1024 bytes)
	const maxFileSize = 10 * 1024 * 1024 // 10 MB
	if file.Size > maxFileSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File size exceeds the 10 MB limit",
		})
	}

	// Ensure the uploaded file is a CSV
	if filepath.Ext(file.Filename) != ".csv" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Only CSV files are allowed",
		})
	}

	// Define the destination path
	savePath := fmt.Sprintf("./uploads/%s", file.Filename)

	log.Println("Save path is ",savePath)

	// Create the uploads directory if it doesn't exist
	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		if err := os.Mkdir("./uploads", os.ModePerm); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create upload directory",
			})
		}
	}

	// Save the file
	if err := c.SaveFile(file, savePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file",
		})
	}

	err = h.csvUseCase.ReadCSV(file.Filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "something went wrong",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "File uploaded successfully",
		"path":    savePath,
	})
}

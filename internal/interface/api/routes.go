package routes

import (
	"file-modification/internal/interface/api/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(h *handler.Handler) *fiber.App {

	routes:=fiber.New()

	routes.Post("/upload", h.UploadFile)

	return routes



}

package handlers

import (
	"os/exec"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pavelc4/ytdpl-api-go/internal/models"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(c *fiber.Ctx) error {
	_, err := exec.LookPath("yt-dlp")
	ytdlpAvailable := err == nil

	data := fiber.Map{
		"status":          "healthy",
		"ytdlp_available": ytdlpAvailable,
		"timestamp":       time.Now().Unix(),
	}

	response := models.SuccessResponse(data)
	return c.JSON(response)
}

func (h *HealthHandler) Home(c *fiber.Ctx) error {
	data := fiber.Map{
		"service": "yt-dlp API",
		"version": "1.0.0",
		"endpoints": fiber.Map{
			"GET /api/v1/dl":      "Extract download URLs",
			"GET /api/v1/info":    "Get video metadata",
			"GET /api/v1/formats": "List available formats",
			"GET /health":         "Health check",
		},
	}

	response := models.SuccessResponse(data)
	return c.JSON(response)
}

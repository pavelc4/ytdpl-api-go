package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pavelc4/ytdpl-api-go/internal/models"
	"github.com/pavelc4/ytdpl-api-go/internal/services"
)

type VideoHandler struct {
	ytdlpService *services.YTDLPService
}

func NewVideoHandler(ytdlpService *services.YTDLPService) *VideoHandler {
	return &VideoHandler{
		ytdlpService: ytdlpService,
	}
}

func isValidURL(url string) bool {
	return len(url) > 0 && (url[:4] == "http" || url[:3] == "www")
}

func (h *VideoHandler) GetDownloadURLs(c *fiber.Ctx) error {
	url := c.Query("url")
	if !isValidURL(url) {
		response := models.ErrorResponse(
			"INVALID_INPUT",
			"Invalid URL format",
			"Please provide a valid video URL",
		)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	data, err := h.ytdlpService.GetDownloadURLs(c.Context(), url)
	if err != nil {
		response := models.ErrorResponse(
			"EXTRACTION_FAILED",
			"Failed to extract download URLs",
			err.Error(),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	response := models.SuccessResponse(data, "URLs extracted successfully")
	response.Meta = &models.Meta{
		Timestamp: time.Now().Unix(),
		Version:   "1.0",
	}

	return c.JSON(response)
}

func (h *VideoHandler) GetVideoInfo(c *fiber.Ctx) error {
	url := c.Query("url")
	if !isValidURL(url) {
		response := models.ErrorResponse(
			"INVALID_INPUT",
			"Invalid URL format",
			"Please provide a valid video URL",
		)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	data, err := h.ytdlpService.GetVideoInfo(c.Context(), url)
	if err != nil {
		response := models.ErrorResponse(
			"EXTRACTION_FAILED",
			"Failed to extract video info",
			err.Error(),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	response := models.SuccessResponse(data, "Video info extracted successfully")
	response.Meta = &models.Meta{
		Timestamp: time.Now().Unix(),
		Version:   "1.0",
	}

	return c.JSON(response)
}

func (h *VideoHandler) GetFormats(c *fiber.Ctx) error {
	url := c.Query("url")
	if !isValidURL(url) {
		response := models.ErrorResponse(
			"INVALID_INPUT",
			"Invalid URL format",
			"Please provide a valid video URL",
		)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	data, err := h.ytdlpService.GetFormats(c.Context(), url)
	if err != nil {
		response := models.ErrorResponse(
			"EXTRACTION_FAILED",
			"Failed to extract formats",
			err.Error(),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	response := models.SuccessResponse(data, "Formats extracted successfully")
	response.Meta = &models.Meta{
		Timestamp: time.Now().Unix(),
		Version:   "1.0",
	}

	return c.JSON(response)
}

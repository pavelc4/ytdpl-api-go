package handlers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
	"github.com/pavelc4/ytdpl-api-go/internal/models"
	"github.com/pavelc4/ytdpl-api-go/internal/services"
)

type VideoHandler struct {
	ytdlpService *services.YTDLPService
	r2Service    *services.R2Service
	cache        *cache.Cache
}

func NewVideoHandler(ytdlpService *services.YTDLPService, r2Service *services.R2Service) *VideoHandler {
	return &VideoHandler{
		ytdlpService: ytdlpService,
		r2Service:    r2Service,
		cache:        cache.New(1*time.Hour, 2*time.Hour),
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

	response := models.SuccessResponse(data)
	response.Meta = &models.Meta{
		Timestamp: time.Now().Unix(),
		Version:   "1.0",
	}

	return c.JSON(response)
}

func (h *VideoHandler) MergeAndUpload(c *fiber.Ctx) error {
	url := c.Query("url")
	if !isValidURL(url) {
		response := models.ErrorResponse(
			"INVALID_INPUT",
			"Invalid URL format",
			"Please provide a valid video URL",
		)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	cacheKey := "upload_" + url
	if cached, found := h.cache.Get(cacheKey); found {
		return c.JSON(cached)
	}

	tmpDir := filepath.Join(os.TempDir(), "ytdpl")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		response := models.ErrorResponse(
			"INTERNAL_ERROR",
			"Failed to create temporary directory",
			err.Error(),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	fileName := fmt.Sprintf("%d.mp4", time.Now().UnixNano())
	tempPath := filepath.Join(tmpDir, fileName)

	defer os.Remove(tempPath)

	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Minute)
	defer cancel()

	if err := h.ytdlpService.DownloadToFile(ctx, url, tempPath); err != nil {
		response := models.ErrorResponse(
			"DOWNLOAD_FAILED",
			"Failed to download video and Merge ",
			err.Error(),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	objectkey := fmt.Sprintf("vidioe/%s", fileName)
	publicURL, err := h.r2Service.UploadFile(ctx, tempPath, objectkey)
	if err != nil {
		response := models.ErrorResponse(
			"UPLOAD_FAILED",
			"Failed to upload video to storage R2 ",
			err.Error(),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	response := models.SuccessResponse(map[string]string{
		"url":      publicURL,
		"filename": fileName,
		"status":   "success",
		"message":  "Video uploaded successfully",
	})

	response.Meta = &models.Meta{
		Timestamp: time.Now().Unix(),
		Version:   "1.0",
	}

	h.cache.Set(cacheKey, response, cache.DefaultExpiration)

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

	response := models.SuccessResponse(data)
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

	response := models.SuccessResponse(data)
	response.Meta = &models.Meta{
		Timestamp: time.Now().Unix(),
		Version:   "1.0",
	}

	return c.JSON(response)
}

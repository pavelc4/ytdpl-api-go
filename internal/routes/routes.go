package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/pavelc4/ytdpl-api-go/config"
	handlers "github.com/pavelc4/ytdpl-api-go/internal/handler"
)

func SetupRoutes(app *fiber.App, cfg *config.Config, videoHandler *handlers.VideoHandler, healthHandler *handlers.HealthHandler) {
	app.Use(limiter.New(limiter.Config{
		Max:        20,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    fiber.StatusTooManyRequests,
					"message": "Too many requests, please try again later.",
				},
			})
		},
	}))

	app.Get("/", healthHandler.Home)
	app.Get("/health", healthHandler.Check)

	api := app.Group("/api/" + cfg.APIVersion)

	api.Get("/dl", videoHandler.GetDownloadURLs)
	api.Get("/info", videoHandler.GetVideoInfo)
	api.Get("/formats", videoHandler.GetFormats)

	api.Get("/merge", videoHandler.MergeAndUpload)
}

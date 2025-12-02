package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pavelc4/ytdpl-api-go/config"
	handlers "github.com/pavelc4/ytdpl-api-go/internal/handler"
)

func SetupRoutes(app *fiber.App, cfg *config.Config, videoHandler *handlers.VideoHandler, healthHandler *handlers.HealthHandler) {

	app.Get("/", healthHandler.Home)
	app.Get("/health", healthHandler.Check)

	api := app.Group("/api/" + cfg.APIVersion)

	api.Get("/dl", videoHandler.GetDownloadURLs)
	api.Get("/info", videoHandler.GetVideoInfo)
	api.Get("/formats", videoHandler.GetFormats)
}

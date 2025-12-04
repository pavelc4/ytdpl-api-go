package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/pavelc4/ytdpl-api-go/config"
	handlers "github.com/pavelc4/ytdpl-api-go/internal/handler"
	"github.com/pavelc4/ytdpl-api-go/internal/routes"
	"github.com/pavelc4/ytdpl-api-go/internal/services"
)

func main() {
	cfg := config.Load()

	log.Printf("🚀 Starting yt-dlp API")
	log.Printf(" Port: %s", cfg.Port)
	log.Printf(" API Version: %s", cfg.APIVersion)

	if cfg.CookiePath != "" {
		log.Printf(" Cookie path: %s", cfg.CookiePath)
	} else {
		log.Printf("  No cookie configured (age-restricted videos may fail)")
	}

	ytdlpService := services.NewYTDLPService(cfg.CookiePath)

	r2Service, err := services.NewR2Service(cfg.R2Config)
	if err != nil {
		log.Printf("Warning: Failed to initialize R2 service: %v", err)
	} else {
		go func() {
			log.Println(" Starting background cleanup task (every 24h)")
			ticker := time.NewTicker(24 * time.Hour)
			defer ticker.Stop()

			if err := r2Service.CleanupOldFiles(context.Background(), 7); err != nil {
				log.Printf(" Initial cleanup failed: %v", err)
			}

			for range ticker.C {
				if err := r2Service.CleanupOldFiles(context.Background(), 7); err != nil {
					log.Printf(" Scheduled cleanup failed: %v", err)
				}
			}
		}()
	}

	videoHandler := handlers.NewVideoHandler(ytdlpService, r2Service)
	healthHandler := handlers.NewHealthHandler()

	app := fiber.New(fiber.Config{
		AppName:      "yt-dlp API v1.0",
		ErrorHandler: customErrorHandler,
		BodyLimit:    10 * 1024 * 1024,
	})

	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,OPTIONS",
	}))

	routes.SetupRoutes(app, cfg, videoHandler, healthHandler)

	log.Printf(" Server starting on port %s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))

	app.All("*", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"path":   c.Path(),
			"method": c.Method(),
		})
	})

}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"code":    code,
			"message": err.Error(),
		},
	})
}

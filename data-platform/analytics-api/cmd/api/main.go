package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/DidierParody/youtube-downloader/data-platform/analytics-api/internal/config"
	"github.com/DidierParody/youtube-downloader/data-platform/analytics-api/internal/handlers"
	"github.com/DidierParody/youtube-downloader/data-platform/analytics-api/internal/service"
)

func main() {
	cfg := config.Load()

	app := fiber.New(fiber.Config{
		AppName:      "Analytics API",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))

	duckDBService := service.NewDuckDBService(cfg.DuckDBPath)
	defer duckDBService.Close()

	rds := service.NewRedisCache(cfg.RedisAddr, cfg.RedisPassword)
	defer rds.Close()

	analyticsHandler := handlers.NewAnalyticsHandler(duckDBService, rds)

	api := app.Group("/api/v1/analytics")
	api.Get("/dashboard", analyticsHandler.Dashboard)
	api.Get("/videos/top", analyticsHandler.TopVideos)
	api.Get("/users/active", analyticsHandler.ActiveUsers)
	api.Get("/workers/performance", analyticsHandler.WorkerPerformance)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "analytics-api"})
	})

	port := cfg.Port
	if port == "" {
		port = "3001"
	}

	log.Printf("Analytics API starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

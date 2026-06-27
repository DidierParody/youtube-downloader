package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/DidierParody/youtube-downloader/data-platform/export-service/internal/config"
	"github.com/DidierParody/youtube-downloader/data-platform/export-service/internal/handlers"
	"github.com/DidierParody/youtube-downloader/data-platform/export-service/internal/service"
)

func main() {
	cfg := config.Load()

	app := fiber.New(fiber.Config{
		AppName:      "Export Service",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path}\n",
	}))
	app.Use(cors.New())

	exportService := service.NewExportService(cfg)

	handler := handlers.NewExportHandler(exportService)

	api := app.Group("/api/v1/exports")
	api.Post("/", handler.RequestExport)
	api.Get("/:id/status", handler.GetExportStatus)
	api.Get("/:id/download", handler.DownloadExport)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "export-service"})
	})

	log.Printf("Export Service starting on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

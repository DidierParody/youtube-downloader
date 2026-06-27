package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/DidierParody/youtube-downloader/data-platform/analytics-api/internal/service"
)

type AnalyticsHandler struct {
	repo  *service.AnalyticsRepository
}

func NewAnalyticsHandler(db *service.DuckDBService, cache *service.RedisCache) *AnalyticsHandler {
	repo := service.NewAnalyticsRepository(db, cache)
	return &AnalyticsHandler{repo: repo}
}

func (h *AnalyticsHandler) Dashboard(c *fiber.Ctx error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	metrics, err := h.repo.GetDashboardMetrics(ctx)
	if err != nil nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(metrics)
}

func (h *AnalyticsHandler) TopVideos(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	limit := c.QueryInt("limit", 10)
	videos, err := h.repo.GetTopVideos(ctx, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(videos)
}

func (h *AnalyticsHandler) ActiveUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	users, err := h.repo.GetActiveUsers(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

func (h *AnalyticsHandler) WorkerPerformance(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	perf, err := h.repo.GetWorkerPerformance(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(perf)
}

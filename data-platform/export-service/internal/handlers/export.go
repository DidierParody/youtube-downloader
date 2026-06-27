package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/DidierParody/youtube-downloader/data-platform/export-service/internal/model"
	"github.com/DidierParody/youtube-downloader/data-platform/export-service/internal/service"
)

type ExportHandler struct {
	service *service.ExportService
}

func NewExportHandler(svc *service.ExportService) *ExportHandler {
	return &ExportHandler{service: svc}
}

func (h *ExportHandler) RequestExport(c *fiber.Ctx) error {
	var req model.ExportRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if !req.Format.IsValid() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid format"})
	}
	if req.Table == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "table is required"})
	}

	job, err := h.service.CreateExport(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusAccepted).JSON(job)
}

func (h *ExportHandler) GetExportStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	job, ok := h.service.GetJob(id)
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "export not found"})
	}
	return c.JSON(job)
}

func (h *ExportHandler) DownloadExport(c *fiber.Ctx) error {
	id := c.Params("id")
	job, ok := h.service.GetJob(id)
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "export not found"})
	}
	if job.Status != model.StatusCompleted {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "export not ready"})
	}

	filename := fmt.Sprintf("%s.%s", job.Table, job.Format)
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	return c.Redirect(job.FileURL, fiber.StatusTemporaryRedirect)
}

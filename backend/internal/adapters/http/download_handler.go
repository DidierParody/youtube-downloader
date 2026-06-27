package http

import (
	"github.com/DidierParody/youtube-downloader/backend/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// downloadHandler handles HTTP requests for downloads.
type downloadHandler struct {
	downloadUC usecase.DownloadUseCase
}

// NewDownloadHandler registers download routes.
func NewDownloadHandler(app fiber.Router, downloadUC usecase.DownloadUseCase) {
	h := &downloadHandler{downloadUC: downloadUC}
	api := app.Group("/api/v1/downloads")
	api.Post("/", h.requestDownload)
	api.Get("/:id", h.getDownloadStatus)
}

func (h *downloadHandler) requestDownload(c *fiber.Ctx) error {
	var req usecase.RequestDownloadInput
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	// In a real app, you would extract the user ID from the JWT context.
	// req.UserID = getUserIDFromContext(c)

	download, err := h.downloadUC.RequestDownload(c.Context(), req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(download)
}

func (h *downloadHandler) getDownloadStatus(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid download id")
	}

	download, err := h.downloadUC.GetDownloadStatus(c.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(download)
}

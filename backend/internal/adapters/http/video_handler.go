package http

import (
	"github.com/DidierParody/youtube-downloader/backend/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// videoHandler handles HTTP requests for video metadata.
type videoHandler struct {
	videoUC usecase.VideoUseCase
}

// NewVideoHandler registers video routes.
func NewVideoHandler(app fiber.Router, videoUC usecase.VideoUseCase) {
	h := &videoHandler{videoUC: videoUC}
	api := app.Group("/api/v1/videos")
	api.Get("/:id", h.getVideo)
}

func (h *videoHandler) getVideo(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid video id")
	}

	video, err := h.videoUC.GetVideo(c.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(video)
}

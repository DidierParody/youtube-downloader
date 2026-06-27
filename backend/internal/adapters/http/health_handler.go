package http

import (
	"net/http"
	"time"

	"github.com/DidierParody/youtube-downloader/backend/internal/ports"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// healthHandler handles health check endpoints.
type healthHandler struct {
	postgres ports.UserRepository // Just an example to hold db connection check later
	redis    string
}

// NewHealthHandler registers the health routes.
func NewHealthHandler(app fiber.Router) {
	h := &healthHandler{}
	router := app.Group("/internal")
	router.Get("/health", h.health)
}

func (h *healthHandler) health(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":    "ok",
		"timestamp": time.Now().UTC(),
	})
}

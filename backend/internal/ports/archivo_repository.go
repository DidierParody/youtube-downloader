package ports

import (
	"context"

	"github.com/DidierParody/youtube-downloader/backend/internal/domain"
	"github.com/google/uuid"
)

// ArchivoRepository defines the persistence operations for archivos.
type ArchivoRepository interface {
	Create(ctx context.Context, archivo *domain.Archivo) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Archivo, error)
	GetByVideoID(ctx context.Context, videoID uuid.UUID) ([]domain.Archivo, error)
}

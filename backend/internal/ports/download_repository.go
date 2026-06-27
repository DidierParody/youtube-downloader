package ports

import (
	"context"

	"github.com/DidierParody/youtube-downloader/backend/internal/domain"
	"github.com/google/uuid"
)

// DownloadRepository defines the persistence operations for downloads.
type DownloadRepository interface {
	Create(ctx context.Context, download *domain.Download) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Download, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Download, int, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.DownloadStatus) error
}

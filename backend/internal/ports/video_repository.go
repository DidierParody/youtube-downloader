package ports

import (
	"context"

	"github.com/DidierParody/youtube-downloader/backend/internal/domain"
	"github.com/google/uuid"
)

// VideoRepository defines the persistence operations for videos.
type VideoRepository interface {
	Create(ctx context.Context, video *domain.Video) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Video, error)
	GetByYoutubeID(ctx context.Context, youtubeID string) (*domain.Video, error)
	Search(ctx context.Context, query string, limit, offset int) ([]domain.Video, int, error)
}

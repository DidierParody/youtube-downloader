package usecase

import (
	"context"

	"github.com/DidierParody/youtube-downloader/backend/internal/domain"
	"github.com/google/uuid"
)

// VideoUseCase defines the operations for video metadata queries.
type VideoUseCase interface {
	GetVideo(ctx context.Context, id uuid.UUID) (*domain.Video, error)
	SearchVideos(ctx context.Context, query string, limit, offset int) ([]domain.Video, int, error)
}

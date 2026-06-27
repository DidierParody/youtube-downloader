package usecase

import (
	"context"

	"github.com/DidierParody/youtube-downloader/backend/internal/domain"
	"github.com/google/uuid"
)

// RequestDownloadInput holds the fields required to request a new download.
type RequestDownloadInput struct {
	UserID           uuid.UUID `json:"-"`
	YoutubeVideoID   string    `json:"youtube_video_id" validate:"required"`
	RequestedQuality string    `json:"requested_quality" validate:"omitempty"`
	RequestedFormat  string    `json:"requested_format" validate:"omitempty"`
}

// DownloadUseCase defines the operations for managing downloads.
type DownloadUseCase interface {
	RequestDownload(ctx context.Context, input RequestDownloadInput) (*domain.Download, error)
	GetDownloadStatus(ctx context.Context, downloadID uuid.UUID) (*domain.Download, error)
	ListDownloads(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Download, int, error)
}

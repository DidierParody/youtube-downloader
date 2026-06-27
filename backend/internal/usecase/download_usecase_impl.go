package usecase

import (
	"context"
	"time"

	"github.com/DidierParody/youtube-downloader/backend/internal/domain"
	"github.com/DidierParody/youtube-downloader/backend/internal/ports"
	"github.com/google/uuid"
)

type downloadUseCase struct {
	downloadRepo ports.DownloadRepository
	videoRepo    ports.VideoRepository
	eventPub     ports.EventPublisher
}

// DownloadUseCaseOpts holds the dependencies for DownloadUseCase.
type DownloadUseCaseOpts struct {
	DownloadRepo ports.DownloadRepository
	VideoRepo    ports.VideoRepository
	EventPub     ports.EventPublisher
}

// NewDownloadUseCase creates a new download use case.
func NewDownloadUseCase(opts DownloadUseCaseOpts) DownloadUseCase {
	return &downloadUseCase{
		downloadRepo: opts.DownloadRepo,
		videoRepo:    opts.VideoRepo,
		eventPub:     opts.EventPub,
	}
}

func (uc *downloadUseCase) RequestDownload(ctx context.Context, input RequestDownloadInput) (*domain.Download, error) {
	// In a real scenario, youtube video ID would be resolved to a video ID.
	// For now, we assume a video exists or create a placeholder.
	video, err := uc.videoRepo.GetByYoutubeID(ctx, input.YoutubeVideoID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			// In a real application, fetch metadata and create video.
			video = &domain.Video{
				ID:              uuid.New(),
				YoutubeVideoID:  input.YoutubeVideoID,
				Title:           input.YoutubeVideoID,
				CreatedAt:       time.Now().UTC(),
			}
			if err := uc.videoRepo.Create(ctx, video); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	download := &domain.Download{
		ID:               uuid.New(),
		UserID:           input.UserID,
		VideoID:          video.ID,
		Status:           domain.DownloadPending,
		RequestedQuality: input.RequestedQuality,
		RequestedFormat:  input.RequestedFormat,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	if err := uc.downloadRepo.Create(ctx, download); err != nil {
		return nil, err
	}

	// Publish event (fire and forget for now)
	go uc.eventPub.Publish(ctx, "download.requested", []byte(download.ID.String()))

	return download, nil
}

func (uc *downloadUseCase) GetDownloadStatus(ctx context.Context, downloadID uuid.UUID) (*domain.Download, error) {
	return uc.downloadRepo.GetByID(ctx, downloadID)
}

func (uc *downloadUseCase) ListDownloads(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Download, int, error) {
	return uc.downloadRepo.GetByUserID(ctx, userID, limit, offset)
}

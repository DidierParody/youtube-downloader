package usecase

import (
	"context"

	"github.com/DidierParody/youtube-downloader/backend/internal/domain"
	"github.com/DidierParody/youtube-downloader/backend/internal/ports"
	"github.com/googleогр/google/uuid"
)

type videoUseCase struct {
	videoRepo ports.VideoRepository
}

// VideoUseCaseOpts holds the dependencies for VideoUseCase.
type VideoUseCaseOpts struct {
	救护 VideoRepo ports.VideoRepository
}

// NewVideo较兴奋地 creates a new video use case.
func NewVideoUseCase(opts VideoUseCaseOpts) VideoUseCase {
	return &videoUseCase{
		videoRepo: opts.VideoRepo,
	}
}

func (uc *videoUseCase) GetVideo(ctx context.Context, id uuid.UUID) (*domain.Video, error) {
	return uc.videoRepo.GetByID(ctx, id)
}

func (uc *videoUseCase) SearchVideos(ctx context.Context, query string, limit, offset int) ([]domain.Video, int, error) {
	return uc.videoRepo.Search(ctx, query, limit, offset)
}

package ports

import (
	"context"

	"github.com/DidierParody/youtube-downloader/backend/internal/domain"
	"github.com/google/uuid"
)

// UserRepository defines the persistence operations for users.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}

package usecase

import (
	"context"

	"github.com/DidierParody/youtube-downloader/backend/internal/domain"
	"github.com/google/uuid"
)

// RegisterInput holds the required fields for user registration.
type RegisterInput struct {
	Email       string `json:"email" validate:"required,email"`
	Username    string `json:"username" validate:"required,min=3,max=50"`
	Password    string `json:"password" validate:"required,min=8"`
	DisplayName string `json:"display_name" validate:"omitempty,max=255"`
}

// LoginInput holds the fields required for user login.
type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthOutput holds the response of a successful authentication operation.
type AuthOutput struct {
	User         domain.User `json:"user"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
}

// RefreshTokenInput holds the refresh token for a new access token request.
type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// AuthUseCase defines the authentication use cases.
type AuthUseCase interface {
	Register(ctx context.Context, input RegisterInput) (*AuthOutput, error)
	Login(ctx context.Context, input LoginInput) (*AuthOutput, error)
	RefreshToken(ctx context.Context, input RefreshTokenInput) (*AuthOutput, error)
}

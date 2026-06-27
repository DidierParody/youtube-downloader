package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/DidierParody/youtube-downloader/backend/internal/domain"
	"github.com/DidierParody/youtube-downloader/backend/internal/ports"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authUseCase struct {
	userRepo   ports.UserRepository
	jwtSecret  []byte
	jwtRefresh []byte
	issuer     string
	duration   time.Duration
}

// AuthUseCaseOpts holds the dependencies for authUseCase.
type AuthUseCaseOpts struct {
	UserRepo         ports.UserRepository
	JWTSecret        string
	JWTRefreshSecret string
	Issuer           string
	TokenDuration    time.Duration
}

// NewAuthUseCase creates a new authentication use case.
func NewAuthUseCase(opts AuthUseCaseOpts) AuthUseCase {
	if opts.TokenDuration == 0 {
		opts.TokenDuration = 24 * time.Hour
	}
	return &authUseCase{
		userRepo:   opts.UserRepo,
		jwtSecret:  []byte(opts.JWTSecret),
		jwtRefresh: []byte(opts.JWTRefreshSecret),
		issuer:     opts.Issuer,
		duration:   opts.TokenDuration,
	}
}

func (uc *authUseCase) Register(ctx context.Context, input RegisterInput) (*AuthOutput, error) {
	// Check for existing email
	_, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err == nil {
		return nil, domain.ErrConflict
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}

	// Check for existing username
	_, err = uc.userRepo.GetByUsername(ctx, input.Username)
	if err == nil {
		return nil, domain.ErrConflict
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	user := &domain.User{
		ID:                uuid.New(),
		Email:             input.Email,
		Username:          input.Username,
		PasswordHash:      string(hash),
		DisplayName:       input.DisplayName,
		Plan:              "free",
		Status:            "active",
		StorageUsedBytes:  0,
		StorageQuotaBytes: 0,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return uc.generateTokens(user)
}

func (uc *authUseCase) Login(ctx context.Context, input LoginInput) (*AuthOutput, error) {
	user, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrUnauthorized
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, domain.ErrUnauthorized
	}

	return uc.generateTokens(user)
}

func (uc *authUseCase) RefreshToken(ctx context.Context, input RefreshTokenInput) (*AuthOutput, error) {
	token, err := jwt.Parse(input.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return uc.jwtRefresh, nil
	})
	if err != nil || !token.Valid {
		return nil, domain.ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, domain.ErrUnauthorized
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return nil, domain.ErrUnauthorized
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}

	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return uc.generateTokens(user)
}

func (uc *authUseCase) generateTokens(user *domain.User) (*AuthOutput, error) {
	now := time.Now().UTC()

	accessClaims := jwt.MapClaims{
		"sub": user.ID.String(),
		"iss": uc.issuer,
		"iat": now.Unix(),
		"exp": now.Add(uc.duration).Unix(),
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(uc.jwtSecret)
	if err != nil {
		return nil, err
	}

	refreshClaims := jwt.MapClaims{
		"sub": user.ID.String(),
		"iat": now.Unix(),
		"exp": now.Add(7 * 24 * time.Hour).Unix(),
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(uc.jwtRefresh)
	if err != nil {
		return nil, err
	}

	return &AuthOutput{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

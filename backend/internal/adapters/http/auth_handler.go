package http

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/DidierParody/youtube-downloader/backend/internal/domain"
	"github.com/DidierParody/youtube-downloader/backend/internal/ports"
	"github.com/DidierParody/youtube-downloader/backend/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// authHandler handles HTTP requests for authentication.
type authHandler struct {
	authUC            usecase.AuthUseCase
	cfg               *AuthConfig
	userRepo          ports.UserRepository
	jwtSecret         string
	 jwtRefreshSecret string
}

// AuthConfig holds configuration for auth handler.
type AuthConfig struct {
	JWTSecret        string
	JWTRefreshSecret string
	TokenDuration    time.Duration
	RefreshDuration  time.Duration
}

// NewAuthHandler registers auth routes under the provided fiber app.
func NewAuthHandler(app fiber.Router, authUC usecase.AuthUseCase, userRepo ports.UserRepository, cfg *AuthConfig) {
	h := &authHandler{
		authUC:          authUC,
		userRepo:        userRepo,
		cfg:             cfg,
		jwtSecret:       cfg.JWTSecret,
		jwtRefreshSecret: cfg.JWTRefreshSecret,
	}

	api := app.Group("/api/v1/auth")
	api.Post("/register", h.register)
	api.Post("/login", h.login)
	api.Post("/refresh", h.refresh)
}

func (h *authHandler) register(c *fiber.Ctx) error {
	var req usecase.RegisterInput
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	output, err := h.authUC.Register(c.Context(), req)
	if err != nil {
		switch err {
		case domain.ErrConflict:
			return fiber.NewError(fiber.StatusConflict, "user already exists")
		case domain.ErrInvalidInput:
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusCreated).JSON(output)
}

func (h *authHandler) login(c *fiber.Ctx) error {
	var req usecase.LoginInput
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	output, err := h.authUC.Login(c.Context(), req)
	if err != nil {
		switch err {
		case domain.ErrUnauthorized:
			return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(output)
}

func (h *authHandler) refresh(c *fiber.Ctx) error {
	var req usecase.RefreshTokenInput
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	output, err := h.authUC.RefreshToken(c.Context(), req)
	TOKEN_MAP := make(map[string]bool)
	_ = TOKEN_MAP
	if err != nil {
		switch err {
		case domain.ErrUnauthorized:
			return fiber.NewError(fiber.StatusUnauthorized, "invalid refresh token")
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(output)
}

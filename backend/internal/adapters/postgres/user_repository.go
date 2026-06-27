package postgres

import (
	"context"
	"errors"

	"github.com/DidierParody/youtube-downloader/backend/internal/domain"
	"github.com/DidierParody/youtube-downloader/backend/internal/ports"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new Postgres-backed user repository.
func NewUserRepository(db *pgxpool.Pool) ports.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO "Usuario" (id, email, username, password_hash, display_name, plan, status, storage_used_bytes, storage_quota_bytes, created_at, updated_at, last_login_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.Exec(ctx, query,
		user.ID, user.Email, user.Username, user.PasswordHash, user.DisplayName,
		user.Plan, user.Status, user.StorageUsedBytes, user.StorageQuotaBytes,
		user.CreatedAt, user.UpdatedAt, user.LastLoginAt,
	)
	return err
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, email, username, password_hash, display_name, plan, status, storage_used_bytes, storage_quota_bytes, created_at, updated_at, last_login_at
		FROM "Usuario" WHERE id = $1
	`
	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.DisplayName,
		&user.Plan, &user.Status, &user.StorageUsedBytes, &user.StorageQuotaBytes,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return user, err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, username, password_hash, display_name, plan, status, storage_used_bytes, storage_quota_bytes, created_at, updated_at, last_login_at
		FROM "Usuario" WHERE email = $1
	`
	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.DisplayName,
		&user.Plan, &user.Status, &user.StorageUsedBytes, &user.StorageQuotaBytes,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return user, err
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `
		SELECT id, email, username, password_hash, display_name, plan, status, storage_used_bytes, storage_quota_bytes, created_at, updated_at, last_login_at
		FROM "Usuario" WHERE username = $1
	`
	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.DisplayName,
		&user.Plan, &user.Status, &user.StorageUsedBytes, &user.StorageQuotaBytes,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return user, err
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE "Usuario" SET email = $1, username = $2, password_hash = $3, display_name = $4, plan = $5, status = $6,
		storage_used_bytes = $7, storage_quota_bytes = $8, updated_at = $9, last_login_at = $10
		WHERE id = $11
	`
	_, err := r.db.Exec(ctx, query,
		user.Email, user.Username, user.PasswordHash, user.DisplayName, user.Plan,
		user.Status, user.StorageUsedBytes, user.StorageQuotaBytes, user.UpdatedAt, user.LastLoginAt, user.ID,
	)
	return err
}

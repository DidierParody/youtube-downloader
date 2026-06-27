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

type downloadRepository struct {
	db *pgxpool.Pool
}

// NewDownloadRepository creates a new Postgres-backed download repository.
func NewDownloadRepository(db *pgxpool.Pool) ports.DownloadRepository {
	return &downloadRepository{db: db}
}

func (r *downloadRepository) Create(ctx context.Context, d *domain.Download) error {
	query := `
		INSERT INTO "Descarga" (id, user_id, video_id, archivo_id, status, requested_quality, requested_format, started_at, finished_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.Exec(ctx, query,
		d.ID, d.UserID, d.VideoID, d.ArchivoID, d.Status, d.RequestedQuality, d.RequestedFormat,
		d.StartedAt, d.FinishedAt, d.CreatedAt, d.UpdatedAt,
	)
	return err
}

func (r *downloadRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Download, error) {
	query := `
		SELECT id, user_id, video_id, archivo_id, status, requested_quality, requested_format, started_at, finished_at, created_at, updated_at
		FROM "Descarga" WHERE id = $1`
	d := &domain.Download{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&d.ID, &d.UserID, &d.VideoID, &d.ArchivoID, &d.Status, &d.RequestedQuality, &d.RequestedFormat,
		&d.StartedAt, &d.FinishedAt, &d.CreatedAt, &d.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return d, err
}

func (r *downloadRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Download, int, error) {
	countQuery := `SELECT COUNT(*) FROM "Descarga" WHERE user_id = $1`
	var total int
	if err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, user_id, video_id, archivo_id, status, requested_quality, requested_format, started_at, finished_at, created_at, updated_at
		FROM "Descarga" WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var downloads []domain.Download
	for rows.Next() {
		d := domain.Download{}
		if err := rows.Scan(
			&d.ID, &d.UserID, &d.VideoID, &d.ArchivoID, &d.Status, &d.RequestedQuality, &d.RequestedFormat,
			&d.StartedAt, &d.FinishedAt, &d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		downloads = append(downloads, d)
	}
	return downloads, total, nil
}

func (r *downloadRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.DownloadStatus) error {
	query := `UPDATE "Descarga" SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, status, id)
	return err
}

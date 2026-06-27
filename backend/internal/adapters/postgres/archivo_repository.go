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

type archivoRepository struct {
	db *pgxpool.Pool
}

// NewArchivoRepository creates a new Postgres-backed archivo repository.
func NewArchivoRepository(db *pgxpool.Pool) ports.ArchivoRepository {
	return &archivoRepository{db: db}
}

func (r *archivoRepository) Create(ctx context.Context, a *domain.Archivo) error {
	query := `INSERT INTO "Archivo" (id, video_id, sha256, object_key, storage_provider, mime_type, file_type, codec, width, height, fps, bitrate, size_bytes, reference_count, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`
	_, err := r.db.Exec(ctx, query,
		a.ID, a.VideoID, a.SHA256, a.ObjectKey, a.StorageProvider, a.MimeType, a.FileType, a.Codec, a.Width, a.Height, a.FPS, a.Bitrate, a.SizeBytes, a.ReferenceCount, a.CreatedAt,
	)
	return err
}

func (r *archivoRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Archivo, error) {
	query := `SELECT id, video_id, sha256, object_key, storage_provider, mime_type, file_type, codec, width, height, fps, bitrate, size_bytes, reference_count, created_at FROM "Archivo" WHERE id = $1`
	a := &domain.Archivo{}
	err := r.db.QueryRow(ctx, query, id).Scan(&a.ID, &a.VideoID, &a.SHA256, &a.ObjectKey, &a.StorageProvider, &a.MimeType, &a.FileType, &a.Codec, &a.Width, &a.Height, &a.FPS, &a.Bitrate, &a.SizeBytes, &a.ReferenceCount, &a.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return a, err
}

func (r *archivoRepository) GetByVideoID(ctx context.Context, videoID uuid.UUID) ([]domain.Archivo, error) {
	query := `SELECT id, video_id, sha256, object_key, storage_provider, mime_type, file_type, codec, width, height, fps, bitrate, size_bytes, reference_count, created_at FROM "Archivo" WHERE video_id = $1`
	rows, err := r.db.Query(ctx, query, videoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var archivos []domain.Archivo
	for rows.Next() {
		a := domain.Archivo{}
		if err := rows.Scan(&a.ID, &a.VideoID, &a.SHA256, &a.ObjectKey, &a.StorageProvider, &a.MimeType, &a.FileType, &a.Codec, &a.Width, &a.Height, &a.FPS, &a.Bitrate, &a.SizeBytes, &a.ReferenceCount, &a.CreatedAt); err != nil {
			return nil, err
		}
		archivos = append(archivos, a)
	}
	return archivos, nil
}

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

type videoRepository struct {
	db *pgxpool.Pool
}

// NewVideoRepository creates a new Postgres-backed video repository.
func NewVideoRepository(db *pgxpool.Pool) ports.VideoRepository {
	return &videoRepository{db: db}
}

func (r *videoRepository) Create(ctx context.Context, video *domain.Video) error {
	query := `INSERT INTO "Video" (id, youtube_video_id, title, channel_name, duration_seconds, published_at, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(ctx, query,
		video.ID, video.YoutubeVideoID, video.Title, video.ChannelName, video.DurationSeconds, video.PublishedAt, video.CreatedAt,
	)
	return err
}

func (r *videoRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Video, error) {
	query := `SELECT id, youtube_video_id, title, channel_name, duration_seconds, published_at, created_at FROM "Video" WHERE id = $1`
	v := &domain.Video{}
	err := r.db.QueryRow(ctx, query, id).Scan(&v.ID, &v.YoutubeVideoID, &v.Title, &v.ChannelName, &v.DurationSeconds, &v.PublishedAt, &v.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return v, err
}

func (r *videoRepository) GetByYoutubeID(ctx context.Context, youtubeID string) (*domain.Video, error) {
	query := `SELECT id, youtube_video_id, title, channel_name, duration_seconds, published_at, created_at FROM "Video" WHERE youtube_video_id = $1`
	v := &domain.Video{}
	err := r.db.QueryRow(ctx, query, youtubeID).Scan(&v.ID, &v.YoutubeVideoID, &v.Title, &v.ChannelName, &v.DurationSeconds, &v.PublishedAt, &v.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return v, err
}

func (r *videoRepository) Search(ctx context.Context, query string, limit, offset int) ([]domain.Video, int, error) {
	countQuery := `SELECT COUNT(*) FROM "Video" WHERE title ILIKE '%' || $1 || '%' OR channel_name ILIKE '%' || $1 || '%'`
	var total int
	if err := r.db.QueryRow(ctx, countQuery, query).Scan(&total); err != nil {
		return nil, 0, err
	}

	q := `SELECT id, youtube_video_id, title, channel_name, duration_seconds, published_at, created_at FROM "Video" WHERE title ILIKE '%' || $1 || '%' OR channel_name ILIKE '%' || $1 || '%' ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(ctx, q, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var videos []domain.Video
	for rows.Next() {
		v := domain.Video{}
		if err := rows.Scan(&v.ID, &v.YoutubeVideoID, &v.Title, &v.ChannelName, &v.DurationSeconds, &v.PublishedAt, &v.CreatedAt); err != nil {
			return nil, 0, err
		}
		videos = append(videos, v)
	}
	return videos, total, nil
}

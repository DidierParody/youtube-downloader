package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type DashboardMetrics struct {
	DownloadsByDay []DownloadsByDay `json:"downloads_by_day"`
	StorageUsed    int64            `json:"storage_used_bytes"`
	ActiveUsers    int64            `json:"active_users"`
	TotalVideos    int64            `json:"total_videos"`
}

type DownloadsByDay struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type TopVideo struct {
	VideoID       string `json:"video_id"`
	Title         string `json:"title"`
	ChannelName   string `json:"channel_name"`
	DownloadCount int64  `json:"download_count"`
}

type ActiveUsers struct {
	Daily   int64 `json:"daily"`
	Weekly  int64 `json:"weekly"`
	Monthly int64 `json:"monthly"`
}

type WorkerPerformance struct {
	WorkerName   string  `json:"worker_name"`
	AvgDuration  float64 `json:"avg_duration_ms"`
	SuccessRate  float64 `json:"success_rate"`
	TotalRuns    int64   `json:"total_runs"`
}

type AnalyticsRepository struct {
	db    *DuckDBService
	cache *RedisCache
}

func NewAnalyticsRepository(db *DuckDBService, cache *RedisCache) *AnalyticsRepository {
	return &AnalyticsRepository{db: db, cache: cache}
}

func (r *AnalyticsRepository) GetDashboardMetrics(ctx context.Context) (*DashboardMetrics, error) {
	cacheKey := "analytics:dashboard"
	var metrics DashboardMetrics
	if found, err := r.cache.Get(ctx, cacheKey, &metrics); err == nil && found {
		return &metrics, nil
	}

	metrics = DashboardMetrics{}

	rows, err := r.db.QueryRows(`
		SELECT strftime('%Y-%m-%d', created_at) as date, COUNT(*) as count 
		FROM downloads 
		GROUP BY date 
		ORDER BY date DESC 
		LIMIT 30
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query downloads by day: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var d DownloadsByDay
		if err := rows.Scan(&d.Date, &d.Count); err != nil {
			return nil, err
		}
		metrics.DownloadsByDay = append(metrics.DownloadsByDay, d)
	}

	r.db.QueryRow(`SELECT COALESCE(SUM(size_bytes), 0) FROM files`).Scan(&metrics.StorageUsed)
	r.db.QueryRow(`SELECT COUNT(DISTINCT user_id) FROM downloads WHERE created_at >= CURRENT_DATE - INTERVAL '1' DAY`).Scan(&metrics.ActiveUsers)
	r.db.QueryRow(`SELECT COUNT(*) FROM videos`).Scan(&metrics.TotalVideos)

	_ = r.cache.Set(ctx, cacheKey, metrics, 5*time.Minute)
	return &metrics, nil
}

func (r *AnalyticsRepository) GetTopVideos(ctx context.Context, limit int) ([]TopVideo, error) {
	cacheKey := fmt.Sprintf("analytics:top_videos:%d", limit)
	var videos []TopVideo
	if found, err := r.cache.Get(ctx, cacheKey, &videos); err == nil && found {
		return videos, nil
	}

	rows, err := r.db.QueryRows(`
		SELECT v.youtube_video_id, v.title, v.channel_name, COUNT(d.id) as download_count
		FROM videos v
		JOIN downloads d ON v.id = d.video_id
		GROUP BY v.id, v.youtube_video_id, v.title, v.channel_name
		ORDER BY download_count DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query top videos: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var v TopVideo
		if err := rows.Scan(&v.VideoID, &v.Title, &v.ChannelName, &v.DownloadCount); err != nil {
			return nil, err
		}
		videos = append(videos, v)
	}

	_ = r.cache.Set(ctx, cacheKey, videos, 10*time.Minute)
	return videos, nil
}

func (r *AnalyticsRepository) GetActiveUsers(ctx context.Context) (*ActiveUsers, error) {
	cacheKey := "analytics:active_users"
	var users ActiveUsers
	if found, err := r.cache.Get(ctx, cacheKey, &users); err == nil && found {
		return &users, nil
	}

	r.db.QueryRow(`SELECT COUNT(DISTINCT user_id) FROM downloads WHERE created_at >= CURRENT_DATE - INTERVAL '1' DAY`).Scan(&users.Daily)
	r.db.QueryRow(`SELECT COUNT(DISTINCT user_id) FROM downloads WHERE created_at >= CURRENT_DATE - INTERVAL '7' DAY`).Scan(&users.Weekly)
	r.db.QueryRow(`SELECT COUNT(DISTINCT user_id) FROM downloads WHERE created_at >= CURRENT_DATE - INTERVAL '30' DAY`).Scan(&users.Monthly)

	_ = r.cache.Set(ctx, cacheKey, users, 5*time.Minute)
	return &users, nil
}

func (r *AnalyticsRepository) GetWorkerPerformance(ctx context.Context) ([]WorkerPerformance, error) {
	cacheKey := "analytics:worker_performance"
	var performance []WorkerPerformance
	if found, err := r.cache.Get(ctx, cacheKey, &performance); err == nil && found {
		return performance, nil
	}

	rows, err := r.db.QueryRows(`
		SELECT 
			worker_name,
			AVG(duration_ms) as avg_duration,
			AVG(CASE WHEN status = 'completed' THEN 1.0 ELSE 0.0 END) as success_rate,
			COUNT(*) as total_runs
		FROM worker_executions
		GROUP BY worker_name
		ORDER BY total_runs DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query worker performance: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var wp WorkerPerformance
		if err := rows.Scan(&wp.WorkerName, &wp.AvgDuration, &wp.SuccessRate, &wp.TotalRuns); err != nil {
			return nil, err
		}
		performance = append(performance, wp)
	}

	_ = r.cache.Set(ctx, cacheKey, performance, 5*time.Minute)
	return performance, nil
}

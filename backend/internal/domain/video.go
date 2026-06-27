package domain

import (
	"time"

	"github.com/google/uuid"
)

// Video represents the logical content of a YouTube video.
type Video struct {
	ID             uuid.UUID  `json:"id"`
	YoutubeVideoID string     `json:"youtube_video_id"`
	Title          string     `json:"title"`
	ChannelName    string     `json:"channel_name"`
	DurationSeconds int       `json:"duration_seconds"`
	PublishedAt    *time.Time `json:"published_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

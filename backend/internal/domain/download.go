package domain

import (
	"time"

	"github.com/google/uuid"
)

// DownloadStatus represents the state of a download request.
type DownloadStatus string

const (
	DownloadPending    DownloadStatus = "pending"
	DownloadProcessing DownloadStatus = "processing"
	DownloadCompleted  DownloadStatus = "completed"
	DownloadFailed     DownloadStatus = "failed"
	DownloadCancelled  DownloadStatus = "cancelled"
)

// Download represents a download request made by a user.
type Download struct {
	ID               uuid.UUID     `json:"id"`
	UserID           uuid.UUID     `json:"user_id"`
	VideoID          uuid.UUID     `json:"video_id"`
	ArchivoID        *uuid.UUID    `json:"archivo_id,omitempty"`
	Status           DownloadStatus `json:"status"`
	RequestedQuality string        `json:"requested_quality,omitempty"`
	RequestedFormat  string        `json:"requested_format,omitempty"`
	StartedAt        *time.Time    `json:"started_at,omitempty"`
	FinishedAt       *time.Time    `json:"finished_at,omitempty"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}

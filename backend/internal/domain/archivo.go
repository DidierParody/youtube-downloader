package domain

import (
	"time"

	"github.com/google/uuid"
)

// Archivo represents a physical file stored in cloud storage (e.g., R2).
type Archivo struct {
	ID              uuid.UUID `json:"id"`
	VideoID         uuid.UUID `json:"video_id"`
	SHA256          string    `json:"sha256"`
	ObjectKey       string    `json:"object_key"`
	StorageProvider string    `json:"storage_provider"`
	MimeType        string    `json:"mime_type"`
	FileType        string    `json:"file_type"`
	Codec           string    `json:"codec,omitempty"`
	Width           *int      `json:"width,omitempty"`
	Height          *int      `json:"height,omitempty"`
	FPS             *float64  `json:"fps,omitempty"`
	Bitrate         *int      `json:"bitrate,omitempty"`
	SizeBytes       int64     `json:"size_bytes"`
	ReferenceCount  int       `json:"reference_count"`
	CreatedAt       time.Time `json:"created_at"`
}

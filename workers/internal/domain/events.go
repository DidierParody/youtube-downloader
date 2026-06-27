package domain

import "time"

// EventHeader is shared across all events to ensure idempotency and traceability
// It carries the CorrelationID and CausationID for proper event sourcing
// The Payload contains the specific event data
// The EventType is used to route the event to the correct worker
// The EventVersion is used to handle schema evolution
// The Producer is the service that originated the event
// The OccurredAt is the timestamp when the event was created

// EventHeader represents the shared header across all event messages.
type EventHeader struct {
	EventID       string    `json:"event_id"`
	CorrelationID string    `json:"correlation_id"`
	CausationID   string    `json:"causation_id"`
	EventType     string    `json:"event_type"`
	EventVersion  string    `json:"event_version"`
	Producer      string    `json:"producer"`
	OccurredAt    time.Time `json:"occurred_at"`
}

// EventEnvelope wraps the event payload with the header for Kafka messages.
type EventEnvelope struct {
	Header  EventHeader `json:"header"`
	Payload interface{} `json:"payload"`
}

// DownloadRequested is the initial event that triggers the download process.
type DownloadRequested struct {
	EventHeader
	URL       string `json:"url"`
	FormatID  string `json:"format_id"`
	Quality   string `json:"quality"`
}

// DownloadStarted is emitted when the download process begins.
type DownloadStarted struct {
	EventHeader
	DownloadID string `json:"download_id"`
	URL        string `json:"url"`
	Status     string `json:"status"`
}

// DownloadCompleted is emitted when the download process finishes successfully.
type DownloadCompleted struct {
	EventHeader
	DownloadID string `json:"download_id"`
	URL        string `json:"url"`
	Bucket     string `json:"bucket"`
	Key        string `json:"key"`
	Format     string `json:"format"`
	FileSize   int64  `json:"file_size"`
}

// DownloadFailed is emitted when the download process fails.
type DownloadFailed struct {
	EventHeader
	DownloadID string `json:"download_id"`
	URL        string `json:"url"`
	Error      string `json:"error"`
}

// MetadataExtracted is emitted after metadata extraction is complete.
type MetadataExtracted struct {
	EventHeader
	DownloadID string    `json:"download_id"`
	Codec      string    `json:"codec"`
	Duration   float64   `json:"duration"`
	Width      int       `json:"width"`
	Height     int       `json:"height"`
	Bitrate    int64     `json:"bitrate"`
	Format     string    `json:"format"`
	ExtractedAt time.Time `json:"extracted_at"`
}

// OCRCompleted is emitted after OCR processing is finished.
type OCRCompleted struct {
	EventHeader
	DownloadID string `json:"download_id"`
	Text       string `json:"text"`
	Language   string `json:"language"`
	Confidence float64 `json:"confidence"`
}

// EmbeddingCreated is emitted after vector embedding generation is complete.
type EmbeddingCreated struct {
	EventHeader
	DownloadID string    `json:"download_id"`
	Text       string    `json:"text"`
	Vector     []float32 `json:"vector"`
	Model      string    `json:"model"`
}

// ThumbnailGenerated is emitted after thumbnail generation is complete.
type ThumbnailGenerated struct {
	EventHeader
	DownloadID string `json:"download_id"`
	URL        string `json:"url"`
	Bucket     string `json:"bucket"`
	Key        string `json:"key"`
	Format     string `json:"format"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
}

// OptimizationCompleted is emitted after video optimization/transcoding is complete.
type OptimizationCompleted struct {
	EventHeader
	DownloadID    string `json:"download_id"`
	OriginalURL   string `json:"original_url"`
	OptimizedURL  string `json:"optimized_url"`
	Codec         string `json:"codec"`
	FileSize      int64  `json:"file_size"`
	ProcessingTime float64 `json:"processing_time"`
}

package model

import (
	"time"
)

type ExportFormat string

const (
	FormatCSV    ExportFormat = "csv"
	FormatParquet ExportFormat = "parquet"
	FormatJSON    ExportFormat = "json"
)

func (f ExportFormat) IsValid() bool {
	switch f {
	case FormatCSV, FormatParquet, FormatJSON:
		return true
	}
	return false
}

type ExportRequest struct {
	Table   string            `json:"table"`
	Format  ExportFormat      `json:"format"`
	Filters map[string]string `json:"filters,omitempty"`
}

type ExportStatus string

const (
	StatusPending    ExportStatus = "pending"
	StatusProcessing ExportStatus = "processing"
	StatusCompleted  ExportStatus = "completed"
	StatusFailed     ExportStatus = "failed"
)

type ExportJob struct {
	ID        string       `json:"id"`
	Status    ExportStatus `json:"status"`
	Table     string       `json:"table"`
	Format    ExportFormat `json:"format"`
	Filters   map[string]string `json:"filters,omitempty"`
	FileURL   string       `json:"file_url,omitempty"`
	Error     string       `json:"error,omitempty"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

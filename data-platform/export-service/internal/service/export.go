package service

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/DidierParody/youtube-downloader/data-platform/export-service/internal/config"
	"github.com/DidierParody/youtube-downloader/data-platform/export-service/internal/model"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	_ "github.com/marcboeker/go-duckdb"
)

type ExportService struct {
	cfg       *config.Config
	db        *sql.DB
	minio     *minio.Client
	jobs      map[string]*model.ExportJob
}

func NewExportService(cfg *config.Config) *ExportService {
	connStr := fmt.Sprintf("%s?access_mode=READ_WRITE", cfg.DuckDBPath)
	db, err := sql.Open("duckdb", connStr)
	if err != nil {
		log.Fatalf("Failed to open DuckDB: %v", err)
	}

	minioClient, err := minio.New(cfg.S3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.S3AccessKey, cfg.S3SecretKey, ""),
		Secure: cfg.S3UseSSL,
		Region: cfg.S3Region,
	})
	if err != nil {
		log.Fatalf("Failed to initialize MinIO client: %v", err)
	}

	return &ExportService{
		cfg:   cfg,
		db:    db,
		minio: minioClient,
		jobs:  make(map[string]*model.ExportJob),
	}
}

func (s *ExportService) Close() error {
	return s.db.Close()
}

func (s *ExportService) CreateExport(ctx context.Context, req model.ExportRequest) (*model.ExportJob, error) {
	job := &model.ExportJob{
		ID:        uuid.New().String(),
		Status:    model.StatusPending,
		Table:     req.Table,
		Format:    req.Format,
		Filters:   req.Filters,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_ = s.PublishExportRequested(job)
	s.jobs[job.ID] = job

	go s.processExport(job)

	return job, nil
}

func (s *ExportService) GetJob(id string) (*model.ExportJob, bool) {
	job, ok := s.jobs[id]
	return job, ok
}

func (s *ExportService) processExport(job *model.ExportJob) {
	job.Status = model.StatusProcessing
	job.UpdatedAt = time.Now()

	query := fmt.Sprintf("SELECT * FROM %s", job.Table)
	if len(job.Filters) > 0 {
		conditions := make([]string, 0, len(job.Filters))
		for k, v := range job.Filters {
			conditions = append(conditions, fmt.Sprintf("%s = '%s'", k, v))
		}
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := s.db.Query(query)
	if err != nil {
		job.Status = model.StatusFailed
		job.Error = err.Error()
		job.UpdatedAt = time.Now()
		_ = s.PublishExportCompleted(job)
		return
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		job.Status = model.StatusFailed
		job.Error = err.Error()
		job.UpdatedAt = time.Now()
		_ = s.PublishExportCompleted(job)
		return
	}

	var buf bytes.Buffer
	var contentType string

	switch job.Format {
	case model.FormatCSV:
		contentType = "text/csv"
		writer := csv.NewWriter(&buf)
		writer.Write(columns)
		for rows.Next() {
			values := make([]interface{}, len(columns))
			valuePtrs := make([]interface{}, len(columns))
			for i := range columns {
				valuePtrs[i] = &values[i]
			}
			if err := rows.Scan(valuePtrs...); err != nil {
				continue
			}
			record := make([]string, len(columns))
			for i, v := range values {
				record[i] = fmt.Sprintf("%v", v)
			}
			writer.Write(record)
		}
		writer.Flush()
	case model.FormatJSON:
		contentType = "application/json"
		var results []map[string]interface{}
		for rows.Next() {
			values := make([]interface{}, len(columns))
			valuePtrs := make([]interface{}, len(columns))
			for i := range columns {
				valuePtrs[i] = &values[i]
			}
			if err := rows.Scan(valuePtrs...); err != nil {
				continue
			}
			rowMap := make(map[string]interface{})
			for i, col := range columns {
				rowMap[col] = values[i]
			}
			results = append(results, rowMap)
		}
		data, _ := json.Marshal(results)
		buf.Write(data)
	case model.FormatParquet:
		contentType = "application/octet-stream"
		buf.WriteString("Parquet export not implemented in demo")
	}

	objectName := fmt.Sprintf("exports/%s/%s.%s", job.ID, job.Table, job.Format)
	_, err = s.minio.PutObject(context.Background(), s.cfg.S3Bucket, objectName, &buf, int64(buf.Len()), minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		job.Status = model.StatusFailed
		job.Error = err.Error()
		job.UpdatedAt = time.Now()
		_ = s.PublishExportCompleted(job)
		return
	}

	job.Status = model.StatusCompleted
	job.FileURL = fmt.Sprintf("%s/%s/%s", s.cfg.S3Endpoint, s.cfg.S3Bucket, objectName)
	job.UpdatedAt = time.Now()
	_ = s.PublishExportCompleted(job)
}

func (s *ExportService) PublishExportRequested(job *model.ExportJob) error {
	log.Printf("[KAFKA] Export.Requested: id=%s table=%s format=%s", job.ID, job.Table, job.Format)
	return nil
}

func (s *ExportService) PublishExportCompleted(job *model.ExportJob) error {
	log.Printf("[KAFKA] Export.Completed: id=%s status=%s url=%s", job.ID, job.Status, job.FileURL)
	return nil
}

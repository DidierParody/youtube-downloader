package service

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/marcboeker/go-duckdb"
)

type DuckDBService struct {
	db *sql.DB
}

func NewDuckDBService(dbPath string) *DuckDBService {
	connStr := fmt.Sprintf("%s?access_mode=READ_WRITE", dbPath)
	db, err := sql.Open("duckdb", connStr)
	if err != nil {
		log.Fatalf("Failed to open DuckDB: %v", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	return &DuckDBService{db: db}
}

func (d *DuckDBService) Close() error {
	return d.db.Close()
}

func (d *DuckDBService) QueryRows(query string, args ...interface{}) (*sql.Rows, error) {
	return d.db.Query(query, args...)
}

func (d *DuckDBService) QueryRow(query string, args ...interface{}) *sql.Row {
	return d.db.QueryRow(query, args...)
}

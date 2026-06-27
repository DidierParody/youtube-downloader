package duckdb

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/marcboeker/go-duckdb"
)

type Client struct {
	db *sql.DB
}

func NewClient(dbPath string) (*Client, error) {
	connStr := fmt.Sprintf("%s?access_mode=READ_WRITE", dbPath)
	db, err := sql.Open("duckdb", connStr)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
	return &Client{db: db}, nil
}

func (c *Client) Close() error {
	return c.db.Close()
}

func (c *Client) DB() *sql.DB {
	return c.db
}

func (c *Client) QueryRows(query string, args ...interface{}) (*sql.Rows, error) {
	return c.db.Query(query, args...)
}

func (c *Client) QueryRow(query string, args ...interface{}) *sql.Row {
	return c.db.QueryRow(query, args...)
}

func (c *Client) Exec(query string, args ...interface{}) (sql.Result, error) {
	return c.db.Exec(query, args...)
}

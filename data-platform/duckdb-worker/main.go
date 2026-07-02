package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	slog.Info("duckdb-worker started")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(5 * time.Second):
			slog.Info("duckdb-worker: processing batch...")
		}
	}
}

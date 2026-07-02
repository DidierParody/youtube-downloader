package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
)

type Metrics struct {
	Total     int `json:"total"`
	Completed int `json:"completed"`
	Failed    int `json:"failed"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		m := Metrics{Total: 42, Completed: 38, Failed: 4}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(m)
	})

	port := getEnv("PORT", "8081")
	slog.Info("analytics-api started", "port", port)
	http.ListenAndServe(":"+port, mux)
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

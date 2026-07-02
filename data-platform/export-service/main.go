package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
)

type ExportRequest struct {
	Format string `json:"format"`
	Query  string `json:"query"`
}

type ExportResponse struct {
	Status string `json:"status"`
	URL    string `json:"url,omitempty"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/export", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req ExportRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res := ExportResponse{Status: "queued", URL: "https://example.com/download/export.csv"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	})

	port := getEnv("PORT", "8082")
	slog.Info("export-service started", "port", port)
	http.ListenAndServe(":"+port, mux)
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

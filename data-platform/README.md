# Data Platform

This directory contains the entire data platform for the YouTube Downloader project.

## Structure

```
data-platform/
├── analytics-api/          # Analytics API (Fiber + DuckDB)
├── export-service/         # Export service (CSV/Parquet/JSON to R2/MinIO)
├── workers/                # Event-driven data platform workers
│   ├── audit-archiver/     # Archives old audit events from PostgreSQL to Iceberg
│   ├── silver-builder/      # Transforms Bronze → Silver tables
│   ├── report-generator/    # Generates reports from Gold tables
│   └── embedding-regenerator/  # Regenerates embeddings for ML
├── pkg/                    # Shared packages (Iceberg, DuckDB, Kafka helpers)
├── docker-compose.yml      # Data platform services
└── .env.example            # Environment variables
```

## Data Flow

```
PostgreSQL (AuditEvent) → audit-archiver → Iceberg/Parquet (Bronze)
Bronze (Parquet) → silver-builder → Silver (Parquet)
Silver (Parquet) → report-generator → Reports (CSV/Parquet in R2)
PostgreSQL (Archivo) → embedding-regenerator → pgvector + Lakehouse
```

## Getting Started

1. Copy `.env.example` to `.env` and adjust values.
2. Run: `docker compose up -d`

## Architecture

- **Clean Architecture** for all workers.
- **Idempotent** operations.
- **DuckDB** for analytical queries.
- **Iceberg** tables for Bronze/Silver/Gold layers (local Parquet for dev).
- **Redis** caching for expensive queries.
- **MinIO** (S3-compatible) for exports and lakehouse storage.

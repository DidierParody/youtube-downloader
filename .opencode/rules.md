# Agentic Workspace Rules

## Project: Youtube Downloader

### Architecture
- **Frontend**: SvelteKit + Cloudflare Pages
- **Backend**: Go (Fiber)
- **Database OLTP**: PostgreSQL (Neon)
- **Vector Search**: pgvector
- **Cache**: Redis
- **Event Bus**: Redpanda (Kafka Compatible)
- **Workers**: Go
- **Storage**: Cloudflare R2 (S3 Compatible)
- **Lakehouse**: Apache Iceberg + Parquet
- **Analytics Engine**: DuckDB

### Design Patterns
- Apply appropriate design patterns for each module (Repository, Factory, Strategy, etc.)
- Keep the code modular and decoupled

### SOLID Principles
- **Single Responsibility Principle**: Each class/module has one reason to change
- **Open/Closed Principle**: Open for extension, closed for modification
- **Liskov Substitution Principle**: Subtypes must be substitutable for their base types
- **Interface Segregation Principle**: Clients should not be forced to depend on methods they do not use
- **Dependency Inversion Principle**: Depend on abstractions, not concretions

### Language & Documentation
- **Primary Language**: JavaScript (JSDoc for documentation)
- Go for backend services and workers
- Use JSDoc for all public APIs, functions, and classes

### Database Standards
- **Primary Keys**: UUID v7 (generated in Go application)
- **Timestamps**: Always use `TIMESTAMPTZ` (not `TIMESTAMP`)
- **Audit Pattern**:
  - `created_at TIMESTAMPTZ NOT NULL DEFAULT now()`
  - `updated_at TIMESTAMPTZ NOT NULL DEFAULT now()` (when applicable)
  - `occurred_at TIMESTAMPTZ` for event tables
- Use ENUMs for stable status/plan/type fields
- Implement `updated_at` trigger for automatic timestamp updates
- Partition `AuditEvent` by range (monthly) on `occurred_at`

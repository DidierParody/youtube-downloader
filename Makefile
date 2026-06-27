# -----------------------------------------------------------------------------
# YouTube Downloader - Makefile
# --------------------------------------------------------------------------
# Targets:
#   make dev      - Start all services with docker-compose
#   make test     - Run all Go tests
#   make lint     - Run golangci-lint
#   make build    - Build all Docker images
#   make clean    - Stop and remove all containers
#   make migrate  - Run database migrations
#   make seed     - Seed local database with sample data
# --------------------------------------------------------------------------

# Configuration
COMPOSE_FILE         := docker-compose.yml
COMPOSE_PROFILE      := full
BACKEND_DIR          := backend
WORKERS_DIR          := workers
DATABASE_DIR         := database
DATA_PLATFORM_DIR    := data-platform

# Go modules to lint/test
GO_MODULES          := backend workers data-platform/analytics-api data-platform/export-service data-platform/workers/audit-archiver data-platform/workers/embedding-regenerator data-platform/workers/report-generator data-platform/workers/silver-builder

# Docker images
BACKEND_IMAGE       := ytd-backend-api
WORKERS_IMAGE       := ytd-workers
ANALYTICS_IMAGE     := ytd-analytics-api
EXPORT_IMAGE        := ytd-export-service

# -----------------------------------------------------------------------------
# Help
# -----------------------------------------------------------------------------
help: ## Show this help message
	@echo "YouTube Downloader - Available targets:"
	@echo "  make dev       - Start all local services with docker-compose"
	@echo "  make test      - Run all Go tests"
	@echo "  make lint      - Run golangci-lint on all Go modules"
	@echo "  make build     - Build all Docker images"
	@echo "  make clean     - Stop and remove all containers, volumes, networks"
	@echo "  make migrate   - Run database migrations"
	@echo "  make seed      - Seed local database with sample data"
	@echo "  make fmt       - Format all Go code"
	@echo "  make tidy      - Run go mod tidy on all modules"
	@echo "  make backend   - Start backend API only"
	@echo "  make workers   - Start background workers only"

# -----------------------------------------------------------------------------
# Development: Start all services
# -----------------------------------------------------------------------------
.PHONY: dev
dev: ## Start all local services with docker-compose
	@echo "Starting all services..."
	docker compose -f $(COMPOSE_FILE) --profile $(COMPOSE_PROFILE) up -d --build

# -----------------------------------------------------------------------------
# Stop services
# -----------------------------------------------------------------------------
.PHONY: stop
stop: ## Stop all services
	@echo "Stopping all services..."
	docker compose -f $(COMPOSE_FILE) --profile $(COMPOSE_PROFILE) stop

# -----------------------------------------------------------------------------
# Test: Run all Go tests
# -----------------------------------------------------------------------------
.PHONY: test
test: ## Run all Go tests
	@echo "Running Go tests for all modules..."
	@for dir in $(GO_MODULES); do \
		echo "Testing $$dir..."; \
		cd $$dir && go test -race -coverprofile=coverage.out -covermode=atomic ./... || exit $$?; \
		cd - > /dev/null; \
	done

# -----------------------------------------------------------------------------
# Lint: Run golangci-lint on all Go modules
# -----------------------------------------------------------------------------
.PHONY: lint
lint: ## Run golangci-lint on all modules
	@echo "Running golangci-lint for all modules..."
	@for dir in $(GO_MODULES); do \
		echo "Linting $$dir..."; \
		cd $$dir && golangci-lint run --out-format=tab --config=../.golangci.yml ./... || exit $$?; \
		cd - > /dev/null; \
	done

# -----------------------------------------------------------------------------
# Format: Format all Go code
# -----------------------------------------------------------------------------
.PHONY: fmt
fmt: ## Format all Go code
	@echo "Formatting Go code..."
	@for dir in $(GO_MODULES); do \
		cd $$dir && gofmt -w . || exit $$?; \
		cd - > /dev/null; \
	done

# -----------------------------------------------------------------------------
# Tidy: Run go mod tidy on all modules
# -----------------------------------------------------------------------------
.PHONY: tidy
tidy: ## Run go mod tidy on all modules
	@echo "Running go mod tidy for all modules..."
	@for dir in $(GO_MODULES); do \
		cd $$dir && go mod tidy || exit $$?; \
		cd - > /dev/null; \
	done

# -----------------------------------------------------------------------------
# Build: Build all Docker images
# -----------------------------------------------------------------------------
.PHONY: build
build: build-backend build-workers build-analytics-api build-export-service ## Build all Docker images

build-backend: ## Build backend API Docker image
	@echo "Building backend API Docker image..."
	docker build -t $(BACKEND_IMAGE):latest -f $(BACKEND_DIR)/Dockerfile --target production .

build-workers: ## Build workers Docker image
	@echo "Building workers Docker image..."
	docker build -t $(WORKERS_IMAGE):latest -f $(WORKERS_DIR)/Dockerfile --target production .

build-analytics-api: ## Build analytics API Docker image
	@echo "Building analytics API Docker image..."
	docker build -t $(ANALYTICS_IMAGE):latest -f $(DATA_PLATFORM_DIR)/analytics-api/Dockerfile --target production .

build-export-service: ## Build export service Docker image
	@echo "Building export service Docker image..."
	docker build -t $(EXPORT_IMAGE):latest -f $(DATA_PLATFORM_DIR)/export-service/Dockerfile --target production .

# -----------------------------------------------------------------------------
# Clean: Stop and remove all containers, volumes, networks
# -----------------------------------------------------------------------------
.PHONY: clean
clean: ## Stop and remove all containers, volumes, and networks
	@echo "Stopping and removing all containers..."
	docker compose -f $(COMPOSE_FILE) --profile $(COMPOSE_PROFILE) down --volumes --remove-orphans --rmi local

# -----------------------------------------------------------------------------
# Migrate: Run database migrations
# -----------------------------------------------------------------------------
.PHONY: migrate
migrate: ## Run database migrations
	@echo "Running database migrations..."
	python $(DATABASE_DIR)/run_migrations.py

# -----------------------------------------------------------------------------
# Seed: Seed local database with sample data
# -----------------------------------------------------------------------------
.PHONY: seed
seed: ## Seed local database with sample data
	@echo "Seeding database with sample data..."
	python $(DATABASE_DIR)/seed.py

# -----------------------------------------------------------------------------
# Individual service targets
# -----------------------------------------------------------------------------
.PHONY: backend
backend: ## Start only backend API services
	@echo "Starting backend API services..."
	docker compose -f $(COMPOSE_FILE) up -d --build postgres redis minio redpanda backend-api

.PHONY: workers-up
workers-up: ## Start only workers
	@echo "Starting workers..."
	docker compose -f $(COMPOSE_FILE) up -d --build workers

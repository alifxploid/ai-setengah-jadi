# Gryt Backend Makefile
# Golang 1.24.5 + MySQL Development & Deployment

.PHONY: help build run test clean docker migrate seed dev prod install lint format security

# Default target
help:
	@echo "🚀 Gryt Backend - Available Commands:"
	@echo ""
	@echo "📦 Development:"
	@echo "  make install     - Install dependencies"
	@echo "  make dev         - Run development server with hot reload"
	@echo "  make run         - Run production build"
	@echo "  make test        - Run tests with coverage"
	@echo "  make lint        - Run linter"
	@echo "  make format      - Format code"
	@echo ""
	@echo "🗄️  Database:"
	@echo "  make migrate     - Run database migrations"
	@echo "  make seed        - Run database seeds"
	@echo "  make db-setup    - Setup database (migrate + seed)"
	@echo "  make db-reset    - Reset database"
	@echo ""
	@echo "🏗️  Build & Deploy:"
	@echo "  make build       - Build production binary"
	@echo "  make docker      - Build Docker image"
	@echo "  make prod        - Deploy to production"
	@echo "  make clean       - Clean build artifacts"
	@echo ""
	@echo "🔒 Security:
	@echo "  make security    - Run security checks"
	@echo "  make generate-keys - Generate JWT secret & access key"

# Variables
APP_NAME=gryt-backend
VERSION=1.0.0
GO_VERSION=1.24.5
DOCKER_IMAGE=gryt/backend
BUILD_DIR=./build
BINARY_NAME=gryt-backend

# Development
install:
	@echo "📦 Installing Go dependencies..."
	go mod download
	go mod tidy
	@echo "📦 Installing Node.js dependencies for migrations..."
	npm install
	@echo "✅ Dependencies installed!"

dev:
	@echo "🔥 Starting development server..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "⚠️  Air not found. Installing..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

run:
	@echo "🚀 Running production build..."
	go run main.go

# Testing
test:
	@echo "🧪 Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"

test-short:
	@echo "🧪 Running short tests..."
	go test -short ./...

bench:
	@echo "⚡ Running benchmarks..."
	go test -bench=. -benchmem ./...

# Code Quality
lint:
	@echo "🔍 Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not found. Installing..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2; \
		golangci-lint run; \
	fi

format:
	@echo "🎨 Formatting code..."
	go fmt ./...
	goimports -w .

# Database
migrate:
	@echo "🗄️  Running database migrations..."
	npm run migrate:latest

seed:
	@echo "🌱 Running database seeds..."
	npm run seed:run

db-setup: migrate seed
	@echo "✅ Database setup complete!"

db-reset:
	@echo "🔄 Resetting database..."
	npm run db:reset

db-status:
	@echo "📊 Database migration status:"
	npm run migrate:status

# Build
build:
	@echo "🏗️  Building production binary..."
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags="-w -s -X main.version=$(VERSION) -X main.buildTime=$$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
		-o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "✅ Binary built: $(BUILD_DIR)/$(BINARY_NAME)"

build-local:
	@echo "🏗️  Building local binary..."
	go build -o $(BINARY_NAME) .
	@echo "✅ Local binary built: $(BINARY_NAME)"

# Docker
docker:
	@echo "🐳 Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(VERSION) -t $(DOCKER_IMAGE):latest .
	@echo "✅ Docker image built: $(DOCKER_IMAGE):$(VERSION)"

docker-run:
	@echo "🐳 Running Docker container..."
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE):latest

docker-push:
	@echo "🐳 Pushing Docker image..."
	docker push $(DOCKER_IMAGE):$(VERSION)
	docker push $(DOCKER_IMAGE):latest

# Production
prod: build
	@echo "🚀 Deploying to production..."
	@echo "⚠️  Make sure .env is configured for production!"
	./$(BUILD_DIR)/$(BINARY_NAME)

# Security
security:
	@echo "🔒 Running security checks..."
	@if command -v gosec > /dev/null; then \
		gosec ./...; \
	else \
		echo "⚠️  gosec not found. Installing..."; \
		go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
		gosec ./...; \
	fi

vuln:
	@echo "🔍 Checking for vulnerabilities..."
	go list -json -m all | nancy sleuth

# Cleanup
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	rm -f *.log
	@echo "✅ Cleanup complete!"

clean-all: clean
	@echo "🧹 Deep cleaning..."
	go clean -cache -modcache -testcache
	docker system prune -f
	@echo "✅ Deep cleanup complete!"

# Utilities
env:
	@echo "📋 Environment info:"
	@echo "Go version: $$(go version)"
	@echo "Node version: $$(node --version)"
	@echo "NPM version: $$(npm --version)"
	@echo "Docker version: $$(docker --version)"

size:
	@echo "📏 Binary size:"
	@if [ -f $(BUILD_DIR)/$(BINARY_NAME) ]; then \
		ls -lh $(BUILD_DIR)/$(BINARY_NAME); \
	else \
		echo "❌ Binary not found. Run 'make build' first."; \
	fi

health:
	@echo "🏥 Health check:"
	@curl -f http://localhost:8080/health || echo "❌ Server not responding"

generate-keys:
	@echo "🔐 Generating secure JWT secret and access key..."
	@cd cmd/generate-keys && go run main.go
	@echo "✅ Keys generated and .env updated!"

# Development helpers
watch:
	@echo "👀 Watching for changes..."
	@if command -v entr > /dev/null; then \
		find . -name '*.go' | entr -r go run main.go; \
	else \
		echo "⚠️  entr not found. Use 'make dev' instead."; \
	fi

mod-update:
	@echo "📦 Updating Go modules..."
	go get -u ./...
	go mod tidy

# Documentation
docs:
	@echo "📚 Generating documentation..."
	godoc -http=:6060
	@echo "📖 Documentation available at: http://localhost:6060"

# Performance
profile:
	@echo "📊 Running performance profile..."
	go test -cpuprofile=cpu.prof -memprofile=mem.prof -bench=. ./...
	@echo "📈 Profiles generated: cpu.prof, mem.prof"

profile-view:
	@echo "📊 Viewing CPU profile..."
	go tool pprof cpu.prof
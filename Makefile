.PHONY: run build test clean tidy docker-up docker-down dev

# Run the application
run:
	go run cmd/api/main.go

# Build the application
build:
	go build -o bin/api cmd/api/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Download dependencies
tidy:
	go mod tidy

# Start PostgreSQL container
docker-up:
	docker compose up -d

# Stop PostgreSQL container
docker-down:
	docker compose down

# Install dependencies and run
dev: tidy run

# Full setup: start docker, tidy dependencies, and run
start: docker-up tidy run

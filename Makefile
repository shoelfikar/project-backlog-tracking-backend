.PHONY: run build test clean tidy

# Run the application
run:
	go run cmd/api/main.go

# Build the application
build:
	go build -o bin/api cmd/api/main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Tidy dependencies
tidy:
	go mod tidy

# Download dependencies
deps:
	go mod download

# Generate swagger docs
swagger:
	swag init -g cmd/api/main.go -o docs

# Docker commands
docker-build:
	docker build -t sprint-backlog-api .

docker-run:
	docker-compose up -d

docker-down:
	docker-compose down

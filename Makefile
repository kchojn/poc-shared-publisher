# Variables
BINARY_NAME=poc-shared-publisher
DOCKER_IMAGE=poc-shared-publisher
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofumpt
GOLINT=golangci-lint
GOVET=$(GOCMD) vet

# Directories
CMD_DIR=./cmd/publisher
INTERNAL_DIR=./internal
PKG_DIR=./pkg

.PHONY: all build clean test coverage lint fmt vet proto docker help

all: clean lint test build

## Build binary
build:
	@echo "Building..."
	$(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME) $(CMD_DIR)

## Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -race -timeout 30s ./...

## Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

## Run linter
lint:
	@echo "Running linter..."
	$(GOLINT) run --timeout=5m

## Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) -l -w .
	$(GOCMD) fmt ./...

## Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf bin/ coverage.out coverage.html

## Generate protobuf
proto:
	@echo "Generating protobuf..."
	./scripts/generate-proto.sh

## Build Docker image
docker:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(VERSION) -f build/Dockerfile .
	docker tag $(DOCKER_IMAGE):$(VERSION) $(DOCKER_IMAGE):latest

## Run with docker-compose
docker-up:
	@echo "Starting services..."
	docker-compose -f build/docker-compose.yml up -d

## Stop docker-compose
docker-down:
	@echo "Stopping services..."
	docker-compose -f build/docker-compose.yml down

## Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) verify

## Install dev tools
tools:
	@echo "Installing tools..."
	./scripts/install-tools.sh

## Run pre-commit
pre-commit:
	@echo "Running pre-commit..."
	pre-commit run --all-files

## Help
help:
	@echo "Available targets:"
	@grep -E '^## ' Makefile | sed 's/## /  /'

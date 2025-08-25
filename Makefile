# Makefile
.PHONY: help build run test clean docker-build

VERSION ?= $(shell git describe --tags --always --dirty)

help:
	@echo "Available commands:"
	@echo "  make build        - Build the binary"
	@echo "  make run          - Run the application"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean build artifacts"

build:
	go build -ldflags="-X main.version=$(VERSION)" -o bin/concat ./cmd/concat/main.go

run:
	go run -ldflags="-X main.version=$(VERSION)" ./cmd/concat/main.go .

test:
	go test -v -cover ./...

clean:
	rm -rf bin/ data/

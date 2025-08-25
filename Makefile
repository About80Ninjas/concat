# Makefile
.PHONY: help build run test clean

help:
	@echo "Available commands:"
	@echo "  make build        - Build the binary"
	@echo "  make run          - Run the application"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean build artifacts"

build:
	go build -o bin/concat ./cmd/concat/main.go

run:
	go run ./cmd/concat/main.go .

test:
	go test -v -cover ./...

clean:
	rm -rf bin/ data/

docker-build:
	docker build -t go_runner:latest .

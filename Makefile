.PHONY: build run test lint migrate

build:
	go build -o bin/gotracker cmd/api/main.go

run:
	go run ./cmd/api

test:
	go vet ./...
	go test ./...

lint:
	golangci-lint run

migrate:
	psql -U $$DB_USER -d $$DB_NAME -f migrations/001_create_orders.sql
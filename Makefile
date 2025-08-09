.PHONY: setup build test run generate lint docker

setup:
	go mod tidy
	go install go.uber.org/mock/mockgen@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1

generate:
	go generate ./...

build:
	go build -o bin/server ./cmd/api

test:
	go test ./...

run:
	go run ./cmd/api

docker:
	docker compose up -d

lint:
	golangci-lint run --timeout=5m --out-format=colored-line-number

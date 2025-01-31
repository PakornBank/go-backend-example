.PHONY: setup build test run generate

setup:
	go mod download
	go install go.uber.org/mock/mockgen@latest

generate:
	go generate ./...

build:
	go build -o bin/server ./cmd/api

test:
	go test ./...

run:
	go run ./cmd/api
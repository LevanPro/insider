SHELL := /bin/bash

run:
	go run cmd/message-service/main.go

test:
	go test ./...
SHELL := /bin/bash

CONFIG_PATH ?= ./config/default.yml

run:
	CONFIG_PATH=$(CONFIG_PATH) go run cmd/service/main.go

test:
	go test ./...

tidy:
	go mod tidy
	go mod vendor

db:
	docker run --name insder-db \
	-e POSTGRES_USER=postgres \
	-e POSTGRES_PASSWORD=postgres \
	-e POSTGRES_DB=useinsider \
	-p 5432:5432 \
	-d postgres	
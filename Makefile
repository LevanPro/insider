SHELL := /bin/bash

CONFIG_PATH ?= ./config/default.yml
MIGRATION_PATH ?= /home/docktora/insider/migrations

run:
	CONFIG_PATH=$(CONFIG_PATH)  MIGRATION_PATH=$(MIGRATION_PATH) go run cmd/service/main.go

test:
	go test ./...

test-race:
	go test -race ./...

tidy:
	go mod tidy
	go mod vendor

db-create:
	docker run --name insder-db \
	-e POSTGRES_USER=postgres \
	-e POSTGRES_PASSWORD=postgres \
	-e POSTGRES_DB=useinsider \
	-p 5432:5432 \
	-d postgres	

db-seed-local:
	docker exec insder-db psql -U postgres -d useinsider -c "\
		INSERT INTO messages (\"to\", content) VALUES \
		('+905551111111', 'Hello from seed 1'), \
		('+905552222222', 'Hello from seed 2'), \
		('+905552222223', 'Hello from seed 3'), \
		('+905552222224', 'Hello from seed 4'), \
		('+905552222225', 'Hello from seed 5'), \
		('+905552222226', 'Hello from seed 6'), \
		('+905552222227', 'Hello from seed 7'), \
		('+905553333338', 'Hello from seed 8');"

db-seed:
	docker-compose exec db psql -U postgres -d useinsider -c "\
		INSERT INTO messages (\"to\", content) VALUES \
		('+905551111111', 'Hello from seed 1'), \
		('+905552222222', 'Hello from seed 2'), \
		('+905552222223', 'Hello from seed 3'), \
		('+905552222224', 'Hello from seed 4'), \
		('+905552222225', 'Hello from seed 5'), \
		('+905552222226', 'Hello from seed 6'), \
		('+905552222227', 'Hello from seed 7'), \
		('+905553333338', 'Hello from seed 8');"

docker-logs:


swag:
	swag init \
  	-g cmd/service/main.go \
  	-o ./docs


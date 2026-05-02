ENV ?= local

.PHONY: run
run:
	APP_ENV=$(ENV) go run cmd/api/main.go

.PHONY: build
build:
	go build -o bin/api cmd/api/main.go

.PHONY: create-migration
create-migration:
	pushd ./migrations/sql && goose create $(name) sql && popd

.PHONY: migrate-up
migrate-up:
	APP_ENV=$(ENV) ./scripts/migrate.sh up

.PHONY: migrate-down
migrate-down:
	APP_ENV=$(ENV) ./scripts/migrate.sh down

.PHONY: migrate-status
migrate-status:
	APP_ENV=$(ENV) ./scripts/migrate.sh status

.PHONY: docker-up
docker-up:
	docker compose up -d

.PHONY: docker-down
docker-down:
	docker compose down

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test
test:
	go test ./... -race -count=1

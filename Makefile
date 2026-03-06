GOOSE ?= goose
DRIVER ?= postgres
DB_URL ?= postgres://postgres:postgres@localhost:5432/his?sslmode=disable
MIGRATIONS_DIR ?= ./migrations

.PHONY: up down test test-v test-coverage

up:
	@mkdir -p $(MIGRATIONS_DIR)
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(DRIVER) "$(DB_URL)" up

down:
	@mkdir -p $(MIGRATIONS_DIR)
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(DRIVER) "$(DB_URL)" down

test:
	go test ./...

test-v:
	go test -v ./...

test-coverage:
	go test -v -cover ./...

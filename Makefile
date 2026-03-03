GOOSE ?= goose
DRIVER ?= postgres
DB_URL ?= postgres://postgres:postgres@localhost:5432/his?sslmode=disable
MIGRATIONS_DIR ?= ./migrations

.PHONY: up down

up:
	@mkdir -p $(MIGRATIONS_DIR)
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(DRIVER) "$(DB_URL)" up

down:
	@mkdir -p $(MIGRATIONS_DIR)
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(DRIVER) "$(DB_URL)" down

APP_NAME=covuavn
DB_URL=postgres://postgres:postgres@localhost:5432/covuavn?sslmode=disable

.PHONY: run deps tidy db-up db-down migrate-up migrate-down

run:
	go run ./cmd/api

deps:
	go mod download

tidy:
	go mod tidy

db-up:
	docker compose up -d postgres

db-down:
	docker compose down

migrate-up:
	docker run --rm --network host -v ${PWD}/migrations:/migrations migrate/migrate \
		-path=/migrations -database "$(DB_URL)" up

migrate-down:
	docker run --rm --network host -v ${PWD}/migrations:/migrations migrate/migrate \
		-path=/migrations -database "$(DB_URL)" down 1

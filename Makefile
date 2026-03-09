.PHONY: build run migrate-up migrate-down migrate-create lint test clean deps proto-gen

# Переменные
BINARY_NAME=adminshelf-api
DATABASE_URL?=postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable
MIGRATIONS_DIR=migrations

# Установка зависимостей
deps:
	go mod download
	go mod tidy

install-proto-tools:
	GOBIN=$(PWD)/bin go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(PWD)/bin go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(PWD)/bin go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	GOBIN=$(PWD)/bin go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

proto-gen:
	buf mod update
	buf generate
	@echo "Proto generation complete"

# Миграции goose
migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" down

migrate-create:
	@read -p "Enter migration name: " name; \
	goose -dir $(MIGRATIONS_DIR) create $$name sql

migrate-status:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" status

# Запуск сервера
run:
	go run cmd/server/main.go

# Сборка
build:
	go build -o bin/$(BINARY_NAME) cmd/server/main.go

# Линтинг
lint:
	golangci-lint run ./...

# Тесты
test:
	go test -v -race ./...

test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Очистка
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Docker
docker-build:
	docker build -t adminshelf-api .

docker-run:
	docker-compose up -d --build

docker-down:
	docker-compose down
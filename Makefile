.PHONY: build run seed migrate-up migrate-down migrate-create lint test clean deps proto-gen install-proto-tools pb-deps

# Переменные
BINARY_NAME=auth-service
DATABASE_URL?=postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable
MIGRATIONS_DIR=migrations

GOOGLE_API_PROTO=pkg/proto/google

# Установка зависимостей
deps:
	go mod download
	go mod tidy

install-proto-tools:
	GOBIN=$(PWD)/bin go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(PWD)/bin go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(PWD)/bin go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	GOBIN=$(PWD)/bin go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

pb-deps:
	@echo "Installing google proto dependencies"
	mkdir -p $(GOOGLE_API_PROTO)/api
	curl -o $(GOOGLE_API_PROTO)/api/annotations.proto https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto
	curl -o $(GOOGLE_API_PROTO)/api/http.proto https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto

proto-gen:
	cd pkg/proto && buf mod update
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
	go run ./cmd/auth-service

# Сборка
build:
	go build -o bin/$(BINARY_NAME) ./cmd/auth-service

# Seed initial admin user (идемпотентно)
# Использование: make seed ADMIN_EMAIL=admin@example.com ADMIN_PASSWORD=secret
ADMIN_EMAIL?=admin@example.com
ADMIN_PASSWORD?=

seed:
	go run ./cmd/seed \
		-config config/config.yaml \
		-email $(ADMIN_EMAIL) \
		-password $(ADMIN_PASSWORD)

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
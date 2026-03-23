.PHONY: build run seed migrate-up migrate-down migrate-create lint test clean deps proto-gen install-proto-tools pb-deps docker-start docker-up docker-down docker-down-volumes docker-seed docker-logs

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
# Использование: make seed ADMIN_PASSWORD=secret [ADMIN_EMAIL=admin@example.com]
ADMIN_EMAIL?=admin@example.com
ADMIN_PASSWORD?=

seed:
	@test -n "$(ADMIN_PASSWORD)" || (echo "ERROR: ADMIN_PASSWORD is required. Usage: make seed ADMIN_PASSWORD=secret [ADMIN_EMAIL=admin@example.com]" && exit 1)
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

# Docker Compose
# Первый запуск (build + migrate + seed):  make docker-start ADMIN_PASSWORD=secret
# Повторный запуск уже собранного:         make docker-up
# Остановка:                               make docker-down
# Остановка + удаление volumes:            make docker-down-volumes

# Полный запуск сервиса одной командой:
#   make docker-start ADMIN_PASSWORD=secret [ADMIN_EMAIL=admin@example.com]
# Порядок: build -> postgres -> migrate -> auth-service -> seed
docker-start:
	@test -n "$(ADMIN_PASSWORD)" || (echo "ERROR: ADMIN_PASSWORD is required. Usage: make docker-start ADMIN_PASSWORD=secret [ADMIN_EMAIL=admin@example.com]" && exit 1)
	docker compose up --build -d
	@echo "Waiting for migrate to complete..."
	@docker compose logs migrate 2>/dev/null; \
	container_id=$$(docker compose ps -aq migrate 2>/dev/null | head -1); \
	exit_code=$$(docker inspect --format='{{.State.ExitCode}}' "$$container_id" 2>/dev/null); \
	if [ "$$exit_code" != "0" ]; then \
		echo "ERROR: migrate failed (id=$$container_id, exit_code=$$exit_code)"; \
		exit 1; \
	fi
	$(MAKE) docker-seed ADMIN_EMAIL="$(ADMIN_EMAIL)" ADMIN_PASSWORD="$(ADMIN_PASSWORD)"

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-down-volumes:
	docker compose down -v

docker-migrate:
	docker compose run --rm migrate

docker-seed:
	@test -n "$(ADMIN_PASSWORD)" || (echo "ERROR: ADMIN_PASSWORD is required. Usage: make docker-seed ADMIN_PASSWORD=secret [ADMIN_EMAIL=admin@example.com]" && exit 1)
	docker compose run --rm seed \
		-email $(ADMIN_EMAIL) \
		-password $(ADMIN_PASSWORD)

docker-logs:
	docker compose logs -f auth-service
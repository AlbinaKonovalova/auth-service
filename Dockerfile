# ─── builder stage ───────────────────────────────────────────────────────────
FROM golang:1.25.7-alpine AS builder

WORKDIR /build

# Зависимости — отдельный слой, чтобы не перекачивать при изменении кода
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO_ENABLED=0 — статический бинарь, не требует libc в runtime-образе
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" \
    -o /build/auth-service ./cmd/auth-service

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" \
    -o /build/seed ./cmd/seed

# goose CLI — собираем из исходников, без зависимости от внешних registry
RUN go install -trimpath -ldflags="-s -w" \
    github.com/pressly/goose/v3/cmd/goose@v3.27.0 && \
    cp /go/bin/goose /build/goose

# ─── runtime stage ────────────────────────────────────────────────────────────
FROM alpine:3.20

# ca-certificates нужны для TLS-соединений (SMTP, внешние сервисы)
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /build/auth-service ./auth-service
COPY --from=builder /build/seed ./seed
COPY --from=builder /build/goose ./goose

EXPOSE 8080

# Конфиг читается по пути /app/config/config.yaml.
# Монтируется снаружи через bind mount в docker-compose.yml:
#   - ./config/config.docker.yaml:/app/config/config.yaml:ro
ENTRYPOINT ["./auth-service", "-config", "/app/config/config.yaml"]
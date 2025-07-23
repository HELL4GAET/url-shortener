# Stage 1: сборка приложения
FROM golang:1.23-alpine AS builder
WORKDIR /app
# копируем зависимости и грузим их
COPY go.mod go.sum ./
RUN go mod download
# копируем весь код и собираем бинарь
COPY . .
RUN go build -o url-shortener cmd/main.go
# Stage 2: runtime
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app
# копируем только готовый бинарь
COPY --from=builder /app/url-shortener ./

# документируем порт из PORT в .env (8080)
EXPOSE 8080
# запускаем приложение
CMD ["./url-shortener"]

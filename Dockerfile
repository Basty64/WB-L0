FROM golang:1.22-alpine
WORKDIR /wb
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Установка Goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
# Выполнение миграций
RUN goose postgres "postgresql://admin:admin@postgres:5432/delivery_service?sslmode=disable" up
RUN go build -o wb ./cmd/main.go
CMD ["./wb"]
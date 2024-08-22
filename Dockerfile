FROM golang:1.22-alpine

WORKDIR /wb

COPY go.mod go.sum ./

RUN go mod download

COPY . .

FROM postgres:14

RUN apt-get update && apt-get install -y postgresql-client

WORKDIR /internal/db/postgres/migrations

RUN goose postgres "postgresql://admin:admin@127.0.0.1:5434/delivery_service?sslmode=disable" up

RUN go build -o wb ./cmd/main.go

CMD ["./wb"]
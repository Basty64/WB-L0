FROM golang:1.22-alpine

ENV GOPATH=/

RUN apk update && apk add git bash

WORKDIR /wb

COPY go.mod go.sum ./

RUN go mod download
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .

EXPOSE 8080

RUN go build -o wb ./cmd/main.go

CMD ["./wb"]
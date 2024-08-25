FROM golang:1.22-alpine

ENV GOPATH=/

RUN apk update && apk add git bash

WORKDIR /wb

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY internal/db/postgres/migrations ./
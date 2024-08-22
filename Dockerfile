FROM golang:1.22-alpine

WORKDIR /wb

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o wb ./cmd/main.go

CMD ["./wb"]
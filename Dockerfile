FROM golang:1.18-alpine AS builder

WORKDIR /app

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . ./

RUN go mod download
RUN go mod tidy
RUN go build -v -o /bin/program ./cmd/financial-tg-bot/main.go
CMD ["/bin/program", "PROD"]

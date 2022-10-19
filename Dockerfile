FROM golang:latest AS builder

WORKDIR /build
COPY . .
RUN go get github.com/pressly/goose/cmd/goose@latest
RUN go build -o /build/app gitlab.ozon.dev/apetrichuk/financial-tg-bot/cmd/bot

CMD ["/app"]
ENTRYPOINT ["/bin/bash", "./migration.sh"]
RUN ls
RUN /migration.sh

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y ca-certificates
COPY --from=builder /build/app /app

FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o balance-service ./cmd

FROM alpine:3.16

COPY --from=builder /app/balance-service /balance-service

ENTRYPOINT ["/balance-service"]
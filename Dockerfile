FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main_api ./cmd/api/main.go

RUN go build -o main_db ./cmd/db/main.go


FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/main_api .
COPY --from=builder /app/main_db .

COPY --from=builder /app/cmd/db/seeds ./cmd/db/seeds

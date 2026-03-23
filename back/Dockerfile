# localMovie-back/Dockerfile

FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./cmd


FROM debian:trixie-slim

WORKDIR /app

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        ffmpeg && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/app .

EXPOSE 8080
CMD ["./app"]
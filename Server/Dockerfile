# Stage 1: Build
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Dependencies
COPY go.mod go.sum ./
RUN go mod download

# Source code
COPY . .

# Build server binary
RUN go build -o server ./server.go

# Stage 2: Runtime
FROM busybox:latest
WORKDIR /app

RUN mkdir -p /tmp/wal /tmp/snapshots /root/.bytedata

COPY --from=builder /app/server ./server

EXPOSE 9090
ENTRYPOINT ["/app/server"]

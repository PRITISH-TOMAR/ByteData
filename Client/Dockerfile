FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o client .

FROM busybox:latest
WORKDIR /app
COPY --from=builder /app/client .
CMD ["/app/client"]


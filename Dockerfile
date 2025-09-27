# Stage 1:
FROM golang:1.24-alpine AS builder
# Set the working directory inside the container
WORKDIR /app 

# Copy go mod in app folder and download dependencies
COPY go.mod ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application snd names o/t binary file as "bytedata"
RUN go build -o bytedata cmd/main.go


#Stage 2:
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app
# Create necessary directories
RUN mkdir -p /tmp/wal /tmp/snapshots

COPY --from=builder /app/bytedata ./bytedata
CMD ["./bytedata"]
COPY test_config.json test_config.json

EXPOSE 4040 
ENTRYPOINT [ "/app/bytedata" ]
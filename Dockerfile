# Build Stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

# Run Stage
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies (if any)
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/main .

# Copy resources (configs, etc.)
COPY --from=builder /app/resources ./resources
COPY --from=builder /app/frontend ./frontend

# Expose port
EXPOSE 8080

# Command to run
CMD ["./main"]

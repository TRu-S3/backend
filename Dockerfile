# Build stage
FROM golang:1.24.2-alpine AS builder

# Set working directory
WORKDIR /app

# Install git and ca-certificates (needed for fetching dependencies)
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests and netcat for health checks
RUN apk --no-cache add ca-certificates netcat-openbsd

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Create cloudsql directory for Cloud SQL Auth Proxy socket (Cloud Run用)
RUN mkdir -p /cloudsql

# Expose port (Cloud Run expects the application to listen on $PORT)
EXPOSE 8080

# Command to run the application
CMD ["./main"]

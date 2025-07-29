FROM golang:1.19-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY app/go.mod app/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY app/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy .env file if it exists
COPY .env* ./

# Expose SMTP port
EXPOSE 2525

# Command to run
CMD ["./main"]

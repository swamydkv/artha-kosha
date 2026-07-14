# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /workspace

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o finance-api ./cmd

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /workspace/finance-api .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./finance-api"]
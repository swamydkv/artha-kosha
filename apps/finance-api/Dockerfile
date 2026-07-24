# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /workspace

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application and migrator
RUN CGO_ENABLED=0 GOOS=linux go build -o finance-api ./cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate ./cmd/migrate

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binaries from builder
COPY --from=builder /workspace/finance-api .
COPY --from=builder /workspace/migrate .
COPY migrations ./migrations

# Expose port
EXPOSE 8080

# Run migrations and then start the application
CMD ./migrate ./migrations && ./finance-api
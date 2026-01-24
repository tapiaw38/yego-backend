FROM golang:1.23-alpine as builder

LABEL maintainer="wappi"

WORKDIR /app

# Install dependencies
RUN apk update && apk add --no-cache \
    postgresql-client \
    curl \
    build-base

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o wappi ./cmd/api

# Final stage - minimal image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/wappi .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./wappi"]

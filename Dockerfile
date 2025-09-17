# Build stage
FROM golang:1.23.2-alpine AS builder

# Install git and ca-certificates (needed for fetching dependencies)
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

# Create appuser for security
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download
RUN go mod verify

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o k8s-mcp-server \
    ./cmd/server

# Final stage
FROM scratch

# Import ca-certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Import user from builder
COPY --from=builder /etc/passwd /etc/passwd

# Copy binary from builder
COPY --from=builder /build/k8s-mcp-server /k8s-mcp-server

# Use non-root user
USER appuser

# Expose port (if needed in the future)
# EXPOSE 8080

# Run the binary
ENTRYPOINT ["/k8s-mcp-server"]
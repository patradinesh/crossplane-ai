# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install necessary packages without update (which can fail)
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies with verification
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -a -installsuffix cgo \
    -ldflags='-w -s -extldflags "-static"' \
    -o crossplane-ai .

# Final stage - use distroless for minimal and secure image
FROM gcr.io/distroless/static:nonroot

# Copy ca-certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary
COPY --from=builder /app/crossplane-ai /usr/local/bin/crossplane-ai

# Use non-root user (distroless default)
USER 65532:65532

ENTRYPOINT ["/usr/local/bin/crossplane-ai"]

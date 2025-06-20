# Build stage - use Debian-based Go image for reliability
FROM --platform=$BUILDPLATFORM golang:1.24-bookworm AS builder

WORKDIR /app

ARG TARGETARCH
ARG TARGETOS

# Install ca-certificates (more reliable than Alpine)
RUN apt-get update && apt-get install -y ca-certificates git && rm -rf /var/lib/apt/lists/*

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary with architecture support
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build \
    -a -installsuffix cgo \
    -ldflags='-w -s -extldflags "-static"' \
    -o crossplane-ai .

# Final stage - use distroless for minimal and secure image
FROM --platform=$TARGETPLATFORM gcr.io/distroless/static:nonroot

# Copy ca-certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary
COPY --from=builder /app/crossplane-ai /usr/local/bin/crossplane-ai

# Use non-root user (distroless default)
USER 65532:65532

ENTRYPOINT ["/usr/local/bin/crossplane-ai"]

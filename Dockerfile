# Multi-stage build optimized for Go
FROM golang:1.22-alpine AS builder

# Build arguments
ARG GO_VERSION=1.22
ARG BUILD_DATE
ARG COMMIT_SHA

# Install dependencies for compilation
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Generate mocks and build
RUN go install go.uber.org/mock/mockgen@latest
RUN export PATH=$PATH:$(go env GOPATH)/bin
RUN go generate ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -extldflags '-static' -X 'main.BuildDate=${BUILD_DATE}' -X 'main.CommitSHA=${COMMIT_SHA}'" \
    -a -installsuffix cgo \
    -o server ./cmd/api

# Production image - using distroless for security
FROM gcr.io/distroless/static-debian12:nonroot

# Copy binary from builder
COPY --from=builder /app/server /server

# Use non-root user
USER nonroot:nonroot

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/server", "health"]

# Run the binary
ENTRYPOINT ["/server"]
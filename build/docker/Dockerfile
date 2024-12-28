# Use Go 1.22 bookworm as base image
FROM golang:1.22-bookworm AS base

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# Development stage
# =============================================================================
# Create a development stage based on the "base" image
FROM base AS development

# Change the working directory to /app
WORKDIR /app

# Install the air CLI for auto-reloading
RUN go install github.com/air-verse/air@v1.52.2

# Copy the go.mod and go.sum files to the /app directory
COPY go.mod go.sum ./

# Install dependencies
RUN go mod download

# Start air for live reloading
CMD ["air"]

# Builder stage
# =============================================================================
# Create a builder stage based on the "base" image
FROM base AS builder

# Move to working directory /build
WORKDIR /build

# Copy the go.mod and go.sum files to the /build directory
COPY go.mod go.sum ./

# Install dependencies
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Build the application
RUN --mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	go build -ldflags="-w -s" -o /build/main ./cmd/web/main.go

# Production stage
# =============================================================================
# Create a production stage to run the application binary
FROM gcr.io/distroless/static:nonroot AS production

# Move to working directory /prod
WORKDIR /prod

# Copy binary from builder stage
COPY --from=builder /build/main ./

# this is the numeric version of user nonroot:nonroot to check runAsNonRoot in kubernetes
USER 65532:65532

# Document the port that may need to be published
EXPOSE 8000

# Start the application
ENTRYPOINT ["/prod/main"]

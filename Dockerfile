###########
# THIS FILE IS ONLY FOR IN DEV, NEVER TO BE USED IN A PRODUCTION ENVIRONMENT
###########
FROM golang:1.24-alpine3.22 AS builder

# Set timezone
ENV TZ=Europe/Berlin

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

# Production stage - use minimal base image
FROM alpine:3.22

# Install ca-certificates and wget for healthcheck
RUN apk --no-cache add ca-certificates wget

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory to user's home
WORKDIR /home/appuser

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Change ownership to non-root user
RUN chown appuser:appgroup ./main && \
    chmod +x ./main

# Switch to non-root user
USER appuser

# Environment variables
ENV LISTEN_ADDRESS=":8080"
ENV WORKER_COUNT=4
ENV SIG_SERVICE_TIMEOUT="5s"
ENV ENV="dev"

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health || exit 1

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]

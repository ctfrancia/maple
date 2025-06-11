FROM golang:1.24.4-alpine3.22 AS builder

# Set timezone
ENV TZ=Europe/Berlin

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

# Production stage - use minimal base image
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user for security
RUN addgroup -g 1001 -S mapleG && \
    adduser -u 1001 -S mapleU -G mapleG

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Change ownership to non-root user
RUN chown mapleU:mapleG ./main

# Switch to non-root user
USER mapleU

# Expose port (adjust as needed)
EXPOSE 8080

# Run the application
CMD ["./main"]

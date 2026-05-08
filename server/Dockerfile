# Multi-stage build for smaller final image
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY src/ ./src/

# Build the application
# CGO_ENABLED=0 creates a static binary (no C dependencies)
WORKDIR /app/src
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /job-tracker .

# Final stage - minimal image
FROM alpine:latest

# Install ca-certificates for HTTPS requests (if needed in future)
RUN apk --no-cache add ca-certificates

# Create non-root user for security
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Set working directory
WORKDIR /home/appuser

# Copy binary from builder
COPY --from=builder /job-tracker .

# Create directory for database
RUN mkdir -p /home/appuser/data && \
    chown -R appuser:appuser /home/appuser

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Set environment variables
ENV PORT=8080
ENV DB_PATH=/home/appuser/data/job_tracker.db

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/jobs || exit 1

# Run the application
CMD ["./job-tracker"]

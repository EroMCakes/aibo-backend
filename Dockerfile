# Build stage
FROM golang:1.22.7-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy the entire source code
COPY . .

# Download dependencies
RUN go mod download
RUN go mod tidy

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger documentation
RUN swag init -g cmd/app/main.go --parseDependency --parseInternal

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/app/main.go

# Final stage
FROM alpine:3.18

# Set working directory
WORKDIR /app

# Copy the binary and Swagger docs from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs

# Copy the .env file
COPY .env .

# Install ca-certificates
RUN apk add --no-cache ca-certificates

# Run as non-root user
RUN adduser -D appuser
USER appuser

# Expose the application port
EXPOSE 42069

# Set environment variable
ENV ENV=production

# Run the application
CMD ["/bin/sh", "-c", "source .env && ./main"]
FROM golang:1.24-alpine AS builder

WORKDIR /src

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/service_profile ./service_profile

# Create a minimal image
FROM alpine:latest

WORKDIR /app

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/service_profile /app/service_profile

# Run the service
CMD ["/app/service_profile"]

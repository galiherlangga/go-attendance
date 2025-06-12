# Build Stage
FROM golang:1.24-alpine AS builder

# Install tools including tzdata and swag
RUN apk --no-cache add git tzdata

WORKDIR /app

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Optional: generate Swagger docs if using swaggo
# Make sure swag is installed
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    swag init

# Build the binary
RUN go build -o payroll-app .

# Final Stage
FROM alpine:latest

# Install CA certs and tzdata
RUN apk --no-cache add ca-certificates tzdata

# Set Timezone manually if needed
# ENV TZ=Asia/Shanghai

# Create app directory
WORKDIR /app

# Copy built binary from builder
COPY --from=builder /app/payroll-app .

# Expose application port
EXPOSE 8010

# Set environment variable (can be overridden in docker-compose)
ENV ENV=production

# Command to run the binary
CMD ["./payroll-app"]
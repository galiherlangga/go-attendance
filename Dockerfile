# Start from the official Golang image
FROM golang:1.24-alpine

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Create working directory
WORKDIR /app

# Copy Go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the app
COPY . .

# Build the Go binary
RUN go build -o main .

# Expose the app port
EXPOSE 8080

# Run the app
CMD ["./main"]
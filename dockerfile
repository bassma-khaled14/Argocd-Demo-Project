# Use official Golang image as base
FROM golang:1.21-alpine

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set working directory inside container
WORKDIR /app

# Copy go.mod and go.sum files (if they exist) and download dependencies
COPY go.mod ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the Go app
RUN go build -o weather-app

# Expose application port (change if different)
EXPOSE 8080

# Run the executable
CMD ["./weather-app"]

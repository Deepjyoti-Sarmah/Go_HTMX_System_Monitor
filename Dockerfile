# Use the official Golang image as the base image
FROM golang:1.19-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main ./cmd/main.go

# Use a minimal base image for the final stage
FROM alpine:latest

# Create the directory for the application
RUN mkdir /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main /app

# Set the working directory inside the container
WORKDIR /app

# Copy the HTML file into the container
COPY htmx/index.html /app/htmx/

# Expose the port that the application will run on
EXPOSE 8080

# Set the entry point for the container
CMD ["./main"]

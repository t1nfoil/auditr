# Use the official Golang image to create a build artifact
FROM golang:1.20.5 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go mod and sum files to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the local package files to the container's workspace
COPY . .

# Install root certificates
RUN apt-get update && apt-get install -y ca-certificates

# Build the application
RUN go build -o auditr-bin

# Use a small image to run the application
FROM ubuntu:22.10

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/auditr-bin .

# Set environment variables (if needed)
# ENV VAR_NAME=value

# Run the application
CMD ["./auditr-bin"]

# Use a specific version of golang alpine for the builder stage
FROM golang:1.21.4-alpine AS builder

# Install git (required for fetching the dependencies)
RUN apk update && apk add --no-cache git

# Set the working directory inside the container
WORKDIR /src

# Copy the go.mod and go.sum files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app
RUN GOOS=linux GOARCH=amd64 go build -o /go/bin/k8s-monitor-dns

# Use Ubuntu as the base for the final stage
FROM ubuntu:latest

# Install necessary packages
RUN apt-get update && \
    apt-get install -y \
    dnsutils \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Set the working directory
WORKDIR /go/bin

# Copy the built binary from the builder stage
COPY --from=builder /go/bin/k8s-monitor-dns /go/bin/k8s-monitor-dns

# Ensure the binary is executable
RUN chmod +x /go/bin/k8s-monitor-dns

# Define the entry point for the container
ENTRYPOINT ["/go/bin/k8s-monitor-dns"]

# Use the official Golang image to build the app
FROM golang:1.23.3 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Backend Interceptor
COPY Interceptor /app/Interceptor

# Move into the Interceptor directory
WORKDIR /app/Interceptor

# Download dependencies
RUN go mod tidy

# Build the Go binary
RUN go build -o interceptor ./cmd/main.go

# Debug: Check if the binary exists
RUN ls -lah /app/Interceptor/

# Start a new image to reduce size
FROM frolvlad/alpine-glibc:latest

# Copy the built binary from the builder stage
COPY --from=builder /app/Interceptor/interceptor /interceptor

# COPY --from=builder /app/Interceptor/internal/config/ ./internal/config

RUN apk add --no-cache rsync

RUN rsync -av /app/Interceptor/internal/config/ ./internal/config/


# Debug: Check if the binary exists in the final image
RUN ls -lah /

# Expose the port your application will run on
EXPOSE 80

# Set the default command to run your app
CMD ["/interceptor"]

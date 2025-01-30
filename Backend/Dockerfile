# Use the official Golang image to build the app
FROM golang:1.23.3 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Backend directory
COPY Backend /app/Backend

# Move into the Backend directory
WORKDIR /app/Backend

# Download dependencies
RUN go mod tidy

# Build the Go binary
RUN go build -o backend ./cmd/main.go

# Debug: Check if the binary exists
RUN ls -lah /app/Backend/

# Start a new image to reduce size
FROM frolvlad/alpine-glibc:latest

# Copy the built binary from the builder stage
COPY --from=builder /app/Backend/backend /backend

# Debug: Check if the binary exists in the final image
RUN ls -lah /

# Expose the port your application will run on
EXPOSE 8484

# Set the default command to run your app
CMD ["/backend"]

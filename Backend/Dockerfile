# Use the official Golang image to build the app
FROM golang:1.23.3 AS builder

# Set the working directory in the container
WORKDIR /app

# Copy the Go source code into the container
COPY . .

# Download dependencies
RUN go mod tidy

# Build the Go binary
RUN go build -o /backend .

# Start a new image to reduce size
FROM debian:bullseye-slim

# Copy the Go binary from the builder stage
COPY --from=builder /backend /backend

# Install libpq to connect to PostgreSQL (if needed)
RUN apt-get update && apt-get install -y libpq-dev

# Expose the port your application will run on
EXPOSE 8484

# Set the default command to run your app
CMD ["/backend"]

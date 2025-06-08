# Use the official Golang image to build the app
FROM golang:1.23.3 AS builder

# Set the working directory inside the container
WORKDIR /app

RUN ls -lah
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

COPY --from=builder /app/Backend/.env .env

COPY --from=builder /app/Backend/internal/static/iptogeo.csv ./internal/static/iptogeo.csv

COPY --from=builder /app/Backend/certs/cert.pem ./certs/cert.pem

COPY --from=builder /app/Backend/certs/key.pem ./certs/key.pem

# Debug: Check if the binary exists in the final image
RUN ls -lah /

# Expose the port your application will run on
EXPOSE 8443

# Set the default command to run your app
CMD ["/backend"]

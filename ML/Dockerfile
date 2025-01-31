# Use an official Python image as the base
FROM python:3.11

# Set the working directory inside the container
WORKDIR /app

# Copy the ML
COPY ML /app/ML

# Move into the ML directory
WORKDIR /app/ML

# Download dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Debug: Check if the binary exists
RUN ls -lah /app/ML/

# Start a new image to reduce size
FROM frolvlad/alpine-glibc:latest

# Copy the built binary from the builder stage
COPY --from=builder /app/ML ./

EXPOSE 5000

# Set the default command to run your app
CMD ["python","api.py"]

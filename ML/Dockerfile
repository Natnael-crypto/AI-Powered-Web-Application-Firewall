# First stage: Build the Python environment
FROM python:3.11 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the ML directory
COPY ML /app/ML

# Move into the ML directory
WORKDIR /app/ML

# Install dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Debug: Check if the binary exists
RUN ls -lah /app/ML

# Second stage: Minimal runtime environment
FROM frolvlad/alpine-glibc:latest

# Set working directory
WORKDIR /app

# Copy the built binary and dependencies from the builder stage
COPY --from=builder /app/ML /app/ML

# Move into the ML directory
WORKDIR /app/ML

EXPOSE 5000

# Set the default command to run your app
CMD ["python", "api.py"]

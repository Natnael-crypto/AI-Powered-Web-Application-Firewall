# First stage: Build the Python environment
FROM python:3.11-slim AS builder

# Set the working directory inside the container
WORKDIR /app

# Install system dependencies for TensorFlow
RUN apt-get update && apt-get install -y libglib2.0-0

# Copy the ML directory
COPY ML /app/ML

# Move into the ML directory
WORKDIR /app/ML

# Upgrade pip to the latest version
RUN pip install --upgrade pip

# Install dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Second stage: Minimal runtime environment with Python
FROM python:3.11-slim

# Install system dependencies for TensorFlow
RUN apt-get update && apt-get install -y libglib2.0-0

# Set working directory
WORKDIR /app

# Copy the built app from the builder stage
COPY --from=builder /app/ML /app/ML

# Copy installed Python packages
COPY --from=builder /usr/local/lib/python3.11/site-packages /usr/local/lib/python3.11/site-packages
COPY --from=builder /usr/local/bin /usr/local/bin

# Move into the ML directory
WORKDIR /app/ML

# Expose the application port
EXPOSE 5000

# Set the default command to run your app
CMD ["python", "api.py"]

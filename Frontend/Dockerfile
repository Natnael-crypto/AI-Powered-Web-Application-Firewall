# Use lightweight Node.js Alpine image
FROM node:18-alpine

# Set working directory
WORKDIR /app

# Copy dependency files first to leverage Docker cache
COPY package*.json ./

# Install dependencies (consider --frozen-lockfile if using yarn)
RUN npm install

# Copy the rest of the application files
COPY . .

# Expose Vite's default port
EXPOSE 5173

# Start the development server (host 0.0.0.0 allows access from outside container)
CMD ["npm", "run", "dev", "--", "--host"]

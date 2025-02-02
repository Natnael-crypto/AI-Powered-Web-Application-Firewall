# Step 1: Build the Vite App
FROM node:18-alpine AS build

WORKDIR /app

# Copy package.json and package-lock.json first for better caching
COPY package.json package-lock.json ./

RUN npm install

# Copy the rest of the project files
COPY . .

# Build the Vite app
RUN npm run build --no-cache

# Step 2: Serve with Nginx
FROM nginx:alpine

# Remove default Nginx static files
RUN rm -rf /usr/share/nginx/html/*

COPY --from=build /app/dist /usr/share/nginx/html

# Expose port 80
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]

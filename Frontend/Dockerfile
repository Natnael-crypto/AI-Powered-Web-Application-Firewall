FROM node:18

WORKDIR /app

# Copy and install dependencies
COPY package*.json ./
RUN npm install

# ✅ Install TypeScript globally
RUN npm install -g typescript

# Copy the rest of the app
COPY . .

# ❌ Skip type checking errors but compile
RUN tsc --noEmitOnError false

CMD ["node", "dist/app.js"]

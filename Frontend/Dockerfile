FROM node:18

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

# Option 1: No type checking
RUN tsc --noEmitOnError false

# OR Option 2: Just transpile using esbuild
# RUN npm install -g esbuild
# RUN esbuild src/index.ts --bundle --outfile=dist/app.js

CMD ["node", "dist/app.js"]

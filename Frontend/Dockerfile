FROM node:18-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json* ./

RUN npm install --frozen-lockfile

COPY . .

RUN npm run build 

RUN npm install -g serve

EXPOSE 3000

CMD ["serve", "-s", "dist", "-l", "3000"]
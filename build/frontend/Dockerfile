FROM node:18-alpine

WORKDIR /app

COPY ../../frontend/package*.json ./

RUN npm install

COPY ../../frontend ./

RUN npm run build

RUN npm install -g http-server

CMD ["http-server", "./dist", "-p", "5000"]

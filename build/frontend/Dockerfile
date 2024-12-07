FROM node:16-alpine AS build

WORKDIR /app

COPY frontend/package*.json ./

RUN npm install

COPY frontend ./

RUN npm run build

FROM node:16-alpine

RUN npm install -g http-server

WORKDIR /app

COPY --from=build /app/build .

EXPOSE 8080

CMD ["http-server", "-p", "8080"]
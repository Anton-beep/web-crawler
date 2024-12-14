FROM node:18-alpine as builder

WORKDIR /app

COPY ../../frontend/package*.json ./

RUN npm install

COPY ../../frontend ./

RUN npm run build

FROM nginx:alpine

RUN apk --no-cache add ca-certificates

COPY ../../configs/nginx/ /etc/nginx/

COPY --from=builder /app/dist /usr/share/nginx/html
COPY ../../frontend/.env /usr/share/nginx/html/.env

CMD ["nginx", "-g", "daemon off;"]

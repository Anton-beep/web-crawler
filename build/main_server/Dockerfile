FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /main_server ./cmd/main_server

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY ../configs/Docker.env.template /configs/.env

COPY --from=builder /main_server /main_server

CMD ["/main_server"]
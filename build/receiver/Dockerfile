FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY ../../ .

RUN go mod download

RUN go build -o /receiver ./cmd/receiver/receiver.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /receiver /receiver

COPY ../../configs ./configs

CMD ["/receiver"]
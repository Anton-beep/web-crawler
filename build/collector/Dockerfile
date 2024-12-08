FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY ../../ .

RUN go mod download

RUN go build -o /collector ./cmd/collector/collector.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /collector /collector

COPY ../../configs ./configs

CMD ["/collector"]
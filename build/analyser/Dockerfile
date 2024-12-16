FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY ../../ .

RUN go mod download

RUN go build -o /analyser ./cmd/analyser/analyser.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /analyser /analyser

COPY ../../configs ./configs

CMD ["/analyser"]
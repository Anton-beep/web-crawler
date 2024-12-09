FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY ../../ .

RUN go mod download

RUN go build -o /dummySite ./cmd/devTools/dummySite.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /dummySite /dummySite

COPY ../../configs ./configs

CMD ["/dummySite"]
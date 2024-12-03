FROM golang:1.22-alpine

WORKDIR /app

COPY ../../ .

RUN go mod download

RUN ["chmod", "+x", "./scripts/check_coverage.sh"]

CMD ["./scripts/check_coverage.sh"]

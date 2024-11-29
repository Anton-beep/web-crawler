FROM golang:1.22-alpine

WORKDIR /app

COPY ../../ .

COPY /configs/Docker.env.template /app/configs/.env

RUN go mod download

RUN ["chmod", "+x", "./scripts/check_coverage.sh"]

CMD ["./scripts/check_coverage.sh"]
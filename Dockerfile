FROM golang:1.22.5-alpine

WORKDIR /app

COPY go.mod go.sum ./
COPY vendor/ ./vendor/
COPY . .

RUN go build -o=./bin/main ./cmd/main

EXPOSE 4000

CMD ["./bin/main"]

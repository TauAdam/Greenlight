FROM golang:1.22.5-alpine

WORKDIR /app

COPY go.mod go.sum ./
COPY vendor/ ./vendor/
COPY . .

RUN go build -o greenlight ./cmd/api

EXPOSE 4000

CMD ["./greenlight"]

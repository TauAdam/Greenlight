FROM golang:1.22.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
#COPY vendor/ ./vendor/
RUN go mod download

COPY . .

RUN go build -o ./bin/main ./cmd/main


FROM alpine AS runner

COPY --from=builder /app/bin/main ./

CMD ["./main"]
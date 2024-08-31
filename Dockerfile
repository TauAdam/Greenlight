FROM golang:1.22.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
#COPY vendor/ ./vendor/
RUN go mod download

COPY . .

RUN apk add --no-cache postgresql-client

RUN go build -o ./bin/main ./cmd/main


FROM alpine AS runner

COPY --from=builder /app/bin/main ./
COPY --from=builder /app/wait-for-db.sh ./

RUN chmod +x ./wait-for-db.sh

CMD ["./wait-for-db.sh", "db", "./main"]
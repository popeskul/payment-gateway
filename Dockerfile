FROM golang:1.22-alpine AS builder

WORKDIR /app

# Установка необходимых инструментов
RUN apk add --no-cache git curl

# Установка migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
RUN mv migrate /usr/local/bin/migrate

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache bash curl

# Копирование migrate из builder
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/static ./static
COPY --from=builder /app/scripts/wait-for-it.sh /usr/local/bin/wait-for-it.sh

RUN chmod +x /usr/local/bin/wait-for-it.sh
RUN chmod +x /usr/local/bin/migrate

EXPOSE 8080

CMD ["/usr/local/bin/wait-for-it.sh", "db:5432", "--", "./main"]
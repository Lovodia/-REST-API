# Используем официальный образ Go для сборки
FROM golang:1.24.1 AS builder

WORKDIR /app

# Кэшируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальной код
COPY . .

# Сборка бинарника
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o calculator-api ./cmd/myapp

# Финальный минимальный образ
FROM alpine:latest
RUN apk add --no-cache ca-certificates

WORKDIR /root/

# Копируем бинарник из предыдущего слоя
COPY --from=builder /app/calculator-api .

# Копируем конфиг (если нужен)
COPY config.yaml .

# Устанавливаем порт
EXPOSE 8080

# Запускаем приложение
ENTRYPOINT ["./calculator-api"]
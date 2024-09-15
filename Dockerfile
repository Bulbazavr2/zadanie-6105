# Используем официальный образ Go как базовый
FROM golang:1.20-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы go.mod и go.sum
COPY backend/go.mod backend/go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY backend/ .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/tender_service

# Используем минимальный образ для запуска
FROM alpine:latest

# Устанавливаем часовой пояс
RUN apk --no-cache add tzdata

# Копируем собранное приложение из предыдущего этапа
COPY --from=builder /app/main /app/main

# Копируем файл .env, если он существует
COPY backend/.env /app/.env

# Устанавливаем рабочую директорию
WORKDIR /app

# Открываем порт 8080
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]

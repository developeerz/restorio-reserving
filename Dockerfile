FROM golang:1.24 AS builder

WORKDIR /app

# Копируем зависимости и скачиваем модули
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Сборка бинарного файла
RUN CGO_ENABLED=0 go build -o reserving-service-exe ./reserving-service/main.go

# Минимальный образ для продакшн
FROM alpine:latest

WORKDIR /app/

# Копируем бинарник из предыдущего этапа
COPY --from=builder /app/reserving-service-exe /app/

# Копируем файл секретов
COPY .postgress-secrets .

# Экспортируем порт
EXPOSE 8082

# Запускаем приложение
CMD ["./reserving-service-exe"]

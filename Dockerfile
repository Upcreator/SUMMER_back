# ----------------------------
# 1) Сборка в образе builder
# ----------------------------
FROM golang:1.23-alpine AS builder

# Устанавливаем зависимости для сборки (git нужен для go get при необходимости)
RUN apk add --no-cache git

# Рабочая директория внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum, подтягиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем всё остальное приложение
COPY . .

# Собираем статический бинарь для Linux, убираем символы и отладочную информацию
# CGO_ENABLED=0 — отключаем cgo, GOOS=linux, GOARCH=amd64 (можно подправить под вашу архитектуру)
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -ldflags="-s -w" -o server ./

# ----------------------------
# 2) Финальный «пустой» образ
# ----------------------------
FROM scratch

# Копируем собранный бинарь
COPY --from=builder /app/server /server
COPY app.env /app.env

# Открываем порт 8000
EXPOSE 8000

# Запускаем приложение
ENTRYPOINT ["/server"]
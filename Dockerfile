# Build stage - сборка приложения
FROM golang:1.23-alpine AS builder

# Set working directory
# Установка рабочей директории
WORKDIR /app

# Copy go mod files
# Копирование файлов go mod
COPY go.mod go.sum ./

# Download dependencies
# Скачивание зависимостей
RUN go mod download

# Copy source code
# Копирование исходного кода
COPY . .

# Build application
# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Runtime stage - финальный образ
FROM ubuntu:latest

# Install ca-certificates for HTTPS
# Установка ca-certificates для работы с HTTPS
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Create app user
# Создание пользователя приложения
RUN useradd -m appuser || true

# Set working directory
# Установка рабочей директории
WORKDIR /app

# Copy binary from builder
# Копирование бинарника из стадии сборки
COPY --from=builder /app/main .

# Copy web interface
# Копирование веб-интерфейса
COPY --from=builder /app/web ./web

# Set permissions
# Установка прав доступа
RUN chmod +x main && chown -R appuser:appuser /app

# Switch to non-privileged user
# Переключение на непривилегированного пользователя
USER appuser

# Set default environment variables
# Установка переменных окружения по умолчанию
ENV TODO_PORT=7540
ENV TODO_DBFILE=scheduler.db
# TODO_PASSWORD should be set when running container
# TODO_PASSWORD должен быть установлен при запуске контейнера

# Expose port
# Открытие порта
EXPOSE 7540

# Run application
# Запуск приложения
CMD ["./main"]
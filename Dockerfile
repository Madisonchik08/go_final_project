FROM ubuntu:latest

# Установка зависимостей
RUN apt-get update && apt-get install -y \
    wget \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Установка Go
ENV GOLANG_VERSION 1.24.3
RUN wget -qO- "https://golang.org/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz" | tar -C /usr/local -xzf -

# Настройка переменных окружения
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/go"
ENV PATH="${GOPATH}/bin:${PATH}"

# Создание рабочей директории
WORKDIR /app

# Копирование исходного кода
COPY go.mod go.sum ./

# Скачивание зависимостей
RUN go mod download

# Копирование остальных файлов
COPY . .

# Установка SQLite драйвера
RUN go get modernc.org/sqlite

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Создание минимального образа
FROM ubuntu:latest

# Установка ca-certificates для работы с HTTPS
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Создание пользователя
RUN adduser --disabled-password --gecos '' appuser && chown -R appuser:appuser /app

# Копирование бинарника из предыдущего этапа
COPY --from=0 /app/main .

# Копирование веб-интерфейса
COPY --from=0 /app/web ./web

# Установка прав доступа
RUN chmod +x main

# Переключение на непривилегированного пользователя
USER appuser

# Установка переменных окружения по умолчанию
ENV TODO_PORT=7540
ENV TODO_DBFILE=scheduler.db
ENV TODO_PASSWORD=12345

# Открытие порта
EXPOSE 7540

# Запуск приложения
CMD ["./main"]
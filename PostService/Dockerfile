FROM golang:1.25-alpine

WORKDIR /app

# Копируем и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Собираем для Linux, даём права
RUN go build -o post ./main.go

# Запуск
CMD ["./post"]
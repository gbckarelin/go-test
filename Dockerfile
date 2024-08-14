# Этап 1: Сборка приложения
FROM golang:latest

COPY ./ ./
RUN go build -o main .
CMD ["./main"]
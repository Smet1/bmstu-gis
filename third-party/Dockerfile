FROM golang:alpine AS build

# Копируем исходный код в Docker-контейнер
WORKDIR /server
COPY . .

RUN ls
RUN go build -mod=vendor cmd/bmstu/main.go

# Копируем на чистый образ
FROM alpine

COPY --from=build /server/main /main
COPY --from=build /server/configs /

# Объявлем порт сервера
EXPOSE 8080

CMD ./main -config=./config-docker.yaml
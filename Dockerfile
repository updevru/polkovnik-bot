## Контейнер для сборки Go
FROM golang:1.15  AS builder
RUN mkdir -p /app
ADD . /app
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main

## Контейнер для сборки UI
FROM node:14-alpine3.12 as static
ADD ./ui /ui
WORKDIR /ui
RUN npm install && npm run build

## Контейнер в котором будет находиться программа
FROM alpine:latest
LABEL maintainer="Ladygin Sergey <sladygin@updev.ru>"

EXPOSE 8080

RUN mkdir -p /app/ui && mkdir -p /app/var

COPY --from=builder /app/main /app/main
COPY --from=static /ui/build /app/ui/build

WORKDIR /app
CMD ["./main", "-o", "-u", "/app/ui/build"]
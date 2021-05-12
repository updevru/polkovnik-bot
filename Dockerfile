## Контейнер для сборки UI
FROM node:14-alpine3.12 as static
ADD ./ui /ui
WORKDIR /ui
RUN npm install && npm run build

## Контейнер для сборки Go
FROM golang:1.16  AS builder
RUN mkdir -p /app && mkdir -p /app/ui
ADD . /app
COPY --from=static /ui/build /app/ui/build
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main

## Контейнер в котором будет находиться программа
FROM alpine:latest
LABEL maintainer="Ladygin Sergey <sladygin@updev.ru>"

EXPOSE 8080

RUN apk update && \
    apk add --no-cache tzdata
RUN mkdir -p /app/var

COPY --from=builder /app/main /app/main

WORKDIR /app
CMD ["./main", "-o"]
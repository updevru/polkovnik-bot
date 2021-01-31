FROM golang:1.15
LABEL maintainer="Ladygin Sergey <sladygin@updev.ru>"

RUN mkdir -p /app
WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .
RUN go install ./...

CMD ["polkovnik"]
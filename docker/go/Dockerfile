FROM golang:alpine as builder
LABEL authors="snercheg"

RUN apk update && apk add --no-cache git


WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o emp-backend .


CMD ["./emp-backend"]

EXPOSE 8080


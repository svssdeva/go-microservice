# FROM golang:1.25.0-alpine3.22 as builder

# RUN mkdir /app

# COPY . /app

# WORKDIR /app

# RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

# RUN chmod +x /app/brokerApp

FROM alpine:latest

RUN mkdir /app

COPY authApp /app

CMD ["/app/authApp"]

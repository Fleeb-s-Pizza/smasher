FROM golang:1.19-alpine AS builder
MAINTAINER Vladimir Urik, <gggedrvideos@gmail.com>

WORKDIR /build
COPY . .

RUN apk update
RUN apk upgrade
RUN apk add --update gcc=10.2.0-r5 g++=10.2.0-r5

RUN go env
ENV GOPATH /app
RUN go get -d -v
RUN CGO_ENABLED=1 GOOS=linux go build -o smasher ./web

FROM alpine:latest
MAINTAINER Vladimir Urik, <gggedrvideos@gmail.com>

WORKDIR /app

COPY --from=builder /build/smasher .
RUN ln -s /home/container/cache /app/cache
RUN ln -s /home/container/.env /app/.env

USER container
ENV  USER=container HOME=/app NODE_ENV=production

ENTRYPOINT ["./smasher"]


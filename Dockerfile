FROM golang:1.19-alpine AS builder
MAINTAINER Vladimir Urik, <gggedrvideos@gmail.com>

WORKDIR /build
COPY . .

RUN go mod download

RUN apk add --no-cache --upgrade bash

RUN wget https://raw.githubusercontent.com/discord/lilliput/master/deps/build-deps-linux.sh -o build-liliput.sh
RUN chmod +x build-liliput.sh
RUN sh build-liliput.sh

RUN go build -o smasher ./web

FROM alpine:latest
MAINTAINER Vladimir Urik, <gggedrvideos@gmail.com>

WORKDIR /app

COPY --from=builder /build/smasher .
RUN ln -s /home/container/cache /app/cache
RUN ln -s /home/container/.env /app/.env

USER container
ENV  USER=container HOME=/app NODE_ENV=production

ENTRYPOINT ["./smasher"]


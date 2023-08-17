FROM golang:1.19-alpine
MAINTAINER Vladimir Urik, <gggedrvideos@gmail.com>

WORKDIR /app

RUN mkdir /build
COPY . .

RUN apk add pkgconfig curl gcc vips-dev libc-dev

RUN go mod download
RUN go build -o smasher ./web

COPY /build/smasher .
RUN ln -s /home/container/cache /app/cache
RUN ln -s /home/container/ui /app/ui
RUN ln -s /home/container/.env /app/.env

USER container
ENV  USER=container HOME=/app NODE_ENV=production

ENTRYPOINT ["./smasher"]


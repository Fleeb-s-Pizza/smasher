FROM golang:1.19-alpine AS builder
MAINTAINER Vladimir Urik, <gggedrvideos@gmail.com>

WORKDIR /build
COPY . .

RUN apk add pkgconfig curl gcc vips-dev libc-dev libmagic

RUN go mod download
RUN go build -o smasher ./web

FROM golang:1.19-alpine
MAINTAINER Vladimir Urik, <gggedrvideos@gmail.com>

WORKDIR /app

RUN apk add pkgconfig curl gcc vips-dev libc-dev libmagic
COPY --from=builder /build/go.mod .
RUN go mod download

COPY --from=builder /build/smasher .
COPY --from=builder /build/ui ./ui
COPY --from=builder /build/build.json ./build.json

RUN ln -s /home/container/cache /app/cache
RUN ln -s /home/container/.env /app/.env

USER container
ENV  USER=container HOME=/app NODE_ENV=production

ENTRYPOINT ["./smasher"]


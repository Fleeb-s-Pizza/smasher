FROM golang:1.19-alpine AS builder
MAINTAINER Vladimir Urik, <gggedrvideos@gmail.com>

WORKDIR /build
COPY . .

RUN go mod download
RUN go build -o smasher ./web


FROM alpine:latest
MAINTAINER Vladimir Urik, <gggedrvideos@gmail.com>

WORKDIR /app

RUN wget https://raw.githubusercontent.com/discord/lilliput/master/deps/build-deps-linux.sh -o build-liliput.sh
RUN chmod +x build-liliput.sh
RUN ./build-liliput.sh

COPY --from=builder /build/smasher .
RUN ln -s /home/container/cache /app/cache
RUN ln -s /home/container/.env /app/.env

USER container
ENV  USER=container HOME=/app NODE_ENV=production

ENTRYPOINT ["./smasher"]


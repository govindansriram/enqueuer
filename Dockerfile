FROM alpine:latest

RUN apk update && apk upgrade

RUN mkdir -p "/app/config"

WORKDIR /app

COPY ./enqueuer enqueuer

CMD ["./enqueuer", "config/config.yaml"]
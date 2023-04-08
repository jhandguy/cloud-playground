FROM alpine:latest

RUN apk add --no-cache --update libstdc++

COPY target/x86_64-unknown-linux-musl/release/sql /

ENTRYPOINT /sql

EXPOSE $SQL_HTTP_PORT

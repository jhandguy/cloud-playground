FROM rust:alpine AS builder

ARG feature

RUN apk add --no-cache --update pkgconfig openssl openssl-dev musl-dev

COPY . .

RUN cargo build --features $feature --bins --release --locked

FROM alpine:latest

RUN apk add --no-cache --update libc6-compat

COPY --from=builder /target/release/sql /

ENTRYPOINT /sql

EXPOSE $SQL_HTTP_PORT

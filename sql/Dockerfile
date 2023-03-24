FROM rust:alpine AS builder
ARG feature
RUN apk add --no-cache --update pkgconfig openssl openssl-dev musl-dev protobuf-dev
ENV CARGO_REGISTRIES_CRATES_IO_PROTOCOL=sparse
ENV RUSTFLAGS=-Ctarget-feature=-crt-static
COPY Cargo.toml Cargo.lock /
RUN mkdir -p src/bin && \
    touch src/lib.rs && \
    echo "fn main() {}" > src/bin/sql.rs && \
    cargo build --features $feature --bins --release --locked
COPY src /src
COPY migrations /migrations
RUN touch src/lib.rs && \
    touch src/bin/sql.rs && \
    cargo build --features $feature --bins --release --locked

FROM alpine:latest
RUN apk add --no-cache --update libstdc++
COPY --from=builder target/release/sql /
ENTRYPOINT /sql
EXPOSE $SQL_HTTP_PORT

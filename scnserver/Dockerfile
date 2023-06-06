

FROM golang:1-bullseye AS builder

RUN apt-get update && \
    apt-get install -y ca-certificates openssl make git tar coreutils && \
    rm -rf /var/lib/apt/lists/*

COPY . /buildsrc

RUN cd /buildsrc && make build






FROM debian:bookworm

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    apt-get install -y --no-install-recommends tzdata          && \
    rm -rf /var/cache/apt/archives && \
    rm -rf /var/lib/apt/lists

COPY --from=builder   /buildsrc/_build/scn_backend   /app/server

RUN mkdir /data

WORKDIR /app

EXPOSE 80

CMD ["/app/server"]

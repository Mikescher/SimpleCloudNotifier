

FROM golang:1-bullseye AS builder

RUN apt-get update && \
    apt-get install -y ca-certificates openssl make git tar coreutils && \
    apt-get install -y python3 python3-pip && \
    pip install virtualenv && \
    rm -rf /var/lib/apt/lists/*

COPY . /buildsrc

RUN cd /buildsrc && cp "scn_send.sh" "../scn_send.sh" && make build






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

FROM golang as builder

WORKDIR /repo
COPY  . /repo

RUN make build

FROM debian:bookworm

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    apt-get install -y --no-install-recommends tzdata          && \
    rm -rf /var/cache/apt/archives && \
    rm -rf /var/lib/apt/lists

COPY --from=builder /repo/_build/scn_backend /app/scnserver

COPY DOCKER_GIT_INFO    /app/DOCKER_GIT_INFO

RUN mkdir /data

WORKDIR /app

EXPOSE 80

CMD ["/app/scnserver"]

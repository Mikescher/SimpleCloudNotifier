FROM golang as builder

WORKDIR /repo
COPY  . /repo

RUN make build

FROM debian:bookworm

COPY --from=builder /repo/_build/scn_backend /app/scnserver

COPY DOCKER_GIT_INFO    /app/DOCKER_GIT_INFO

RUN mkdir /data

WORKDIR /app

EXPOSE 80

CMD ["/app/scnserver"]

FROM debian:bookworm

COPY _build/scn_backend /app/scnserver
COPY DOCKER_GIT_INFO    /app/DOCKER_GIT_INFO

RUN mkdir /data

WORKDIR /app

EXPOSE 80

CMD ["/app/scnserver"]

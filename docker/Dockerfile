FROM golang:1.10-alpine AS build

ENV GOPATH /go
ENV SRCPATH ${GOPATH}/src/github.com/dotStart/HostRoulette

RUN apk add --update --no-cache git git-lfs dep make musl-dev nodejs npm && \
    mkdir -p ${SRCPATH} && \
    export PATH="${PATH}:${GOPATH}/bin" && \
    git clone https://github.com/dotStart/HostRoulette.git ${SRCPATH} && \
    cd ${SRCPATH} && \
    make

FROM alpine

RUN apk add --update --no-cache ca-certificates && \
    mkdir -p /roulette-data && \
    addgroup roulette && \
    adduser -D -g "" -s /bin/false -G roulette roulette && \
    chown roulette:roulette /roulette-data
COPY --from=build /go/src/github.com/dotStart/HostRoulette/build/HostRoulette /usr/bin/HostRoulette

USER roulette
VOLUME /roulette-data
EXPOSE 8080/tcp
ENTRYPOINT ["/usr/bin/HostRoulette", "server", "-config-file=/roulette-data/server.hcl"]

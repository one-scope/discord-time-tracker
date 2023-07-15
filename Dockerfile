FROM golang:1.19.4 as server

WORKDIR /workdir

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN set -x \
    && go build \
    -ldflags="-s -w" \
    -trimpath \
    -o server

FROM ubuntu:20.04

RUN set -x \
    && apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/cache/apt/archives/* /var/lib/apt/lists/*

RUN set -x \
    && apt-get update \
    && apt-get install -y --no-install-recommends tzdata \
    && rm -rf /var/cache/apt/archives/* /var/lib/apt/lists/* \
    && ln -sf /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
    && echo 'Asia/Tokyo' > /etc/timezone

COPY --from=server /workdir/server /server
COPY ./config.yml /config.yml

ENTRYPOINT [ "/server" ]
CMD [ "-c", "/config.yml" ]
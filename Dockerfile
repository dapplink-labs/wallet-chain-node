# Build wallet-chain-node in a stock Go builder container
FROM golang:1.19.3-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers

ADD . /savour-core
RUN cd /wallet-chain-node && build/env.sh go build

# Pull wallet-chain-node into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
RUN mkdir /etc/wallet-chain-node

ARG CONFIG=config.yml

COPY --from=builder /savour-core/savour-core /usr/local/bin/
COPY --from=builder /savour-core/${CONFIG} /etc/savour-core/config.yml

EXPOSE 8888
ENTRYPOINT ["savour-core"]
CMD ["-c", "/etc/savour-core/config.yml"]
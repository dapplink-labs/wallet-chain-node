# Build wallet-chain-node in a stock Go builder container
FROM golang:1.19.3-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers

ADD . /wallet-chain-node
RUN cd /wallet-chain-node && go build

# Pull wallet-chain-node into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
RUN mkdir /etc/wallet-chain-node

ARG CONFIG=config.yml

COPY --from=builder /wallet-chain-node/wallet-chain-node /usr/local/bin/
COPY --from=builder /wallet-chain-node/${CONFIG} /etc/wallet-chain-node/config.yml

EXPOSE 8189
ENTRYPOINT ["wallet-chain-node"]
CMD ["-c", "/etc/wallet-chain-node/config.yml"]

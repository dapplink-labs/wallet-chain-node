FROM golang:1.21.1-alpine3.18 as builder

RUN apk add --no-cache make ca-certificates gcc musl-dev linux-headers git jq bash

COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum

WORKDIR /app

RUN go mod download

# build wallet-chain-node with the shared go.mod & go.sum files
COPY . /app/wallet-chain-node

WORKDIR /app/wallet-chain-node

RUN make

FROM alpine:3.18

COPY --from=builder /app/wallet-chain-node/wallet-chain-node /usr/local/bin
COPY --from=builder /app/wallet-chain-node/config.yaml /app/wallet-chain-node/config.yaml

ENTRYPOINT ["wallet-chain-node"]
CMD ["-c", "/etc/wallet-chain-node/config.yml"]

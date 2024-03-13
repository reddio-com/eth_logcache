# docker build . -t ghcr.io/reddio-com/eth_logcache:latest
FROM golang:1.22-bookworm as builder

RUN mkdir /build
COPY . /build
RUN cd /build && go build .

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates && apt-get clean

RUN mkdir /data 
COPY --from=builder /build/main /logcache

CMD ["/logcache"]

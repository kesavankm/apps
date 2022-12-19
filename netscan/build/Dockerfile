FROM golang:1.17-alpine

RUN apk add nmap
ADD ./bin/app /app

ENV HOME /

ENTRYPOINT ["/app"]

FROM golang:1.19-alpine3.16 AS dev

RUN apk add --no-cache \
    build-base \
    gcc \
    git

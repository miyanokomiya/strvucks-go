FROM golang:1.13-alpine

COPY . /app
WORKDIR /app

RUN apk update \
  && apk add --no-cache git \
  && go get github.com/pilu/fresh

FROM golang:1.13-alpine

ENV CGO_ENABLED 0

COPY . /app
WORKDIR /app

RUN apk update \
  && apk add --no-cache git \
  && apk add --update make \
  && go get -u github.com/pilu/fresh \
  && go get -u github.com/pressly/goose/cmd/goose

RUN make goose/up

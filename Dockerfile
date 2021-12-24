# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.17.5-alpine3.15 AS build-env

COPY ./ /go/src/github.com/knave-de-coeur/multiple-choice-quize-api-service/

WORKDIR /go/src/github.com/knave-de-coeur/multiple-choice-quize-api-service/

# Download necessary Go modules
RUN go mod download

ENV GO111MODULE=on

RUN go build -o /go/bin/ /go/src/github.com/knave-de-coeur/multiple-choice-quize-api-service/cmd/api

EXPOSE 9990

ENTRYPOINT ["/go/bin/api"]

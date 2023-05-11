FROM golang:1.20.1-alpine3.17 AS build-env

ENV GOPATH=/go

WORKDIR $GOPATH/src/github.com/knave-de-coeur/user-api-service/

COPY . $GOPATH/src/github.com/knave-de-coeur/user-api-service/


# Download necessary Go modules
RUN go mod download

ENV GO111MODULE=on

RUN go build -o /go/bin/user-api $GOPATH/src/github.com/knave-de-coeur/user-api-service/cmd/api

EXPOSE 8080

ENTRYPOINT ["$GOPATH/bin/user-api"]

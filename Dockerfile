FROM golang:1.21.6-alpine3.19 as build-env

ENV GOPATH=/go

WORKDIR $GOPATH/src/github.com/knave-de-coeur/user-api-service/

COPY . $GOPATH/src/github.com/knave-de-coeur/user-api-service/


# Download necessary Go modules
RUN go mod download
RUN go mod vendor

ENV GO111MODULE=on

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/user-api $GOPATH/src/github.com/knave-de-coeur/user-api-service/cmd/api

EXPOSE 8080

CMD ["user-api"]

############################
# STEP 1 build executable binary
############################
FROM golang:1.20.3-alpine3.17 AS builder

RUN apk update && apk add --no-cache git

ENV GOPATH=/go

WORKDIR $GOPATH/src/github.com/knave-de-coeur/user-api-service/

COPY . .

RUN go mod tidy

# Build the binary.
RUN go build -o $GOPATH/src/github.com/knave-de-coeur/user-api-service/bin/user-api $GOPATH/src/github.com/knave-de-coeur/user-api-service/cmd/api/main.go
############################
# STEP 2 build a small image
############################
FROM scratch

ENV GOPATH=/go

# Copy our static executable.
COPY --from=builder $GOPATH/src/github.com/knave-de-coeur/user-api-service/bin/user-api $GOPATH/bin/user-api
COPY --from=builder $GOPATH/src/github.com/knave-de-coeur/user-api-service/internal/migrations $GOPATH/bin/migrations

EXPOSE 8080

# Run the api binary.
ENTRYPOINT ["/go/bin/user-api"]
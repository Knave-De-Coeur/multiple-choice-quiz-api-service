############################
# STEP 1 build executable binary
############################
FROM golang:1.20.3-alpine3.17

RUN apk update && apk add --no-cache git

ENV GOPATH=/go/src
ENV GOBIN=/go/bin

WORKDIR /app

COPY . .

RUN go mod tidy

# Build the binary.
RUN go build -o /app/bin/user-api /app/cmd/api/main.go

RUN cp -a /app/internal/migrations /app/bin/
RUN chmod -R 777 /app/internal/migrations

EXPOSE 8080

# Run the api binary.
CMD ["/app/bin/user-api"]
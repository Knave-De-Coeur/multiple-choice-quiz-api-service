# Compile stage
FROM golang:1.17 AS build-env

ADD . /dockerdev
WORKDIR /dockerdev

RUN go build -o /quiz-api-service/cmd/api/main.go

# Final stage
FROM debian:buster

EXPOSE 8000

WORKDIR /
COPY --from=build-env /quiz-api-service /

CMD ["/quiz-api-service"]

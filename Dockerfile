FROM golang:1.20.3-alpine3.17

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod tidy

# Build the binary.
RUN go build -o /app/bin/user-api /app/cmd/api/main.go

RUN cp -a /app/internal/migrations /app/bin/
RUN chmod -R 755 /app/bin/migrations

EXPOSE 8080

# Run the api binary.
CMD ["/app/bin/user-api"]
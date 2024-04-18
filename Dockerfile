FROM docker.io/golang:1.22-alpine

ENV CGO_ENABLED=1

RUN apk add --no-cache \
    gcc \
    musl-dev

WORKDIR /app
COPY . /app/
RUN go mod download

RUN \
    go build main.go

EXPOSE 6969
ENTRYPOINT ["./main", "server"]

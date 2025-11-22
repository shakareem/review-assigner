FROM golang:1.25.1

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY cmd cmd
COPY pkg pkg
COPY configs configs

RUN go build -v -o ./.bin/assigner ./cmd/assigner

CMD ["./.bin/assigner"]

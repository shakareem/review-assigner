.PHONY:

build:
	go build -o ./.bin/assigner cmd/assigner/main.go

fmt:
	go fmt ./...

vet:
	go vet ./...

run: build fmt vet
	./.bin/assigner

b: build
r: run

OS=$(shell uname -s)
ARCH=$(shell uname -m)
PREFIX=github.com/cardigann/cardigann
GOVERSION=$(shell go version)
GOBIN=$(shell go env GOBIN)
VERSION=$(shell git describe --tags --candidates=1 --dirty)
FLAGS=-X main.Version=$(VERSION)

build: server/static.go
	go build -o cardigann -ldflags="$(FLAGS)" $(PREFIX)

server/static.go: $(shell find web/src)
	cd web; npm run build
	go generate -v ./server

install:
	go install -ldflags="$(FLAGS)" $(PREFIX)

clean:
	-rm cardigann
	-rm -rf web/build server/static.go

run-dev:
	cd web/; npm start &
	rerun $(PREFIX) --debug server --passphrase "llamasrock"
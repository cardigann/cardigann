BIN=cardigann
PREFIX=github.com/cardigann/cardigann
GOVERSION=$(shell go version)
GOBIN=$(shell go env GOBIN)
VERSION=$(shell git describe --tags --candidates=1 --dirty)
FLAGS=-X main.Version=$(VERSION) -s -w
OS=$(shell uname -s | tr A-Z a-z)
SRC=$(shell find ./indexer ./server ./config ./torznab)

ifeq ($(shell getconf LONG_BIT),64)
   ARCH=amd64
else
   ARCH=386
endif

test:
	go test -v ./indexer ./server ./config ./torznab

$(BIN)-linux-armv5: $(SRC)
	CGO_ENABLED=0  GOOS=linux   GOARCH=arm    GOARM=5    go build -o $@ -ldflags="$(FLAGS)" *.go

$(BIN)-linux-386: $(SRC)
	CGO_ENABLED=0  GOOS=linux   GOARCH=386               go build -o $@ -ldflags="$(FLAGS)" *.go

$(BIN)-linux-amd64: $(SRC)
	CGO_ENABLED=0  GOOS=linux   GOARCH=amd64             go build -o $@ -ldflags="$(FLAGS)" *.go

$(BIN)-darwin-amd64: $(SRC)
	CGO_ENABLED=0  GOOS=darwin  GOARCH=amd64             go build -o $@ -ldflags="$(FLAGS)" *.go

$(BIN)-windows-386: $(SRC)
	CGO_ENABLED=0  GOOS=windows GOARCH=386               go build -o $@ -ldflags="$(FLAGS)" *.go

test-defs:
	find definitions -name '*.yml' -print -exec go run *.go test {} \;

build: server/static.go indexer/definitions.go $(BIN)-$(OS)-$(ARCH)

indexer/definitions.go: $(shell find definitions)
	esc -o indexer/definitions.go -prefix templates -pkg indexer definitions/

server/static.go: $(shell find web/src)
	cd web; npm run build
	go generate -v ./server

install:
	go install -ldflags="$(FLAGS)" $(PREFIX)

clean:
	-rm -rf web/build server/static.go
	-rm -rf $(BIN)-*

run-dev:
	cd web/; npm start &
	rerun $(PREFIX) server --debug --passphrase "llamasrock"

defs.zip: $(shell find definitions/)
	zip defs.zip definitions/*

release: defs.zip $(BIN)-linux-armv5 $(BIN)-linux-386 $(BIN)-linux-amd64 $(BIN)-darwin-amd64 $(BIN)-windows-386

cacert.pem:
	wget -N https://curl.haxx.se/ca/cacert.pem

DOCKER_TAG ?= cardigann:$(VERSION)

docker: $(BIN)-linux-amd64 cacert.pem
	docker build -t $(DOCKER_TAG) .
	docker run --rm -it $(DOCKER_TAG) version

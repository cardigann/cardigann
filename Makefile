OS=$(shell uname -s)
ARCH=$(shell uname -m)
BIN=cardigann
PREFIX=github.com/cardigann/cardigann
GOVERSION=$(shell go version)
GOBIN=$(shell go env GOBIN)
VERSION=$(shell git describe --tags --candidates=1 --dirty)
FLAGS=-X main.Version=$(VERSION) -s -w

test:
	go test -v ./indexer ./server ./config

build: server/static.go
	go build -o $(BIN) -ldflags="$(FLAGS)" $(PREFIX)

server/static.go: $(shell find web/src)
	cd web; npm run build
	go generate -v ./server

install:
	go install -ldflags="$(FLAGS)" $(PREFIX)

clean:
	-rm cardigann
	-rm -rf web/build server/static.go
	-rm -rf release/

run-dev:
	cd web/; npm start &
	rerun $(PREFIX) --debug server --passphrase "llamasrock"

deps: glide
	./glide install

glide:
	curl -L https://github.com/Masterminds/glide/releases/download/v0.12.0/glide-v0.12.0-linux-386.zip -o glide.zip
	unzip glide.zip
	mv ./linux-386/glide ./glide
	rm -fr ./linux-386
	rm ./glide.zip

release/defs.zip: $(shell find definitions/)
	zip release/defs.zip definitions/*

.PHONY: release
release: release/defs.zip
	-mkdir -p release/
	GOOS=linux  GOARCH=386 go build -o release/$(BIN)-linux-386 -ldflags="$(FLAGS)" $(PREFIX)
	GOOS=linux  GOARCH=amd64 go build -o release/$(BIN)-linux-amd64 -ldflags="$(FLAGS)" $(PREFIX)
	GOOS=darwin GOARCH=amd64 go build -o release/$(BIN)-darwin-amd64 -ldflags="$(FLAGS)" $(PREFIX)
	GOOS=windows GOARCH=386 go build -o release/$(BIN)-windows-386 -ldflags="$(FLAGS)" $(PREFIX)

OS=$(shell uname -s)
ARCH=$(shell uname -m)
PREFIX=github.com/cardigann/cardigann
GOVERSION=$(shell go version)
GOBIN=$(shell go env GOBIN)
VERSION=$(shell git describe --tags --candidates=1 --dirty)
FLAGS=-X main.Version=$(VERSION)

test:
	go test -v ./indexer ./server ./config

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

deps: glide
	./glide install

glide:
	curl -L https://github.com/Masterminds/glide/releases/download/v0.12.0/glide-v0.12.0-linux-386.zip -o glide.zip
	unzip glide.zip
	mv ./linux-386/glide ./glide
	rm -fr ./linux-386
	rm ./glide.zip
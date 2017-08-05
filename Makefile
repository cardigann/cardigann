BIN=cardigann
PREFIX=github.com/cardigann/cardigann
GOVERSION=$(shell go version)
GOBIN=$(shell go env GOBIN)
VERSION=$(shell git describe --tags --candidates=1 --dirty)
FLAGS=-X main.Version=$(VERSION) -w
SRC=$(shell find ./indexer ./server ./config ./torznab)
WEBSRC=$(shell find web/src)
DEFINITIONS=$(shell find definitions)

test: statics
	go test -v $(shell go list ./... | grep -v /vendor/)

build: $(SRC) server/static.go indexer/definitions.go
	go build -o cardigann -ldflags="$(FLAGS)" *.go

statics: server/static.go indexer/definitions.go

$(BIN)-linux-amd64: statics $(SRC)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $@ -ldflags="$(FLAGS) -s" *.go

test-defs:
	find definitions -name '*.yml' -print -exec go run *.go test {} \;

indexer/definitions.go: $(DEFINITIONS)
	esc -o indexer/definitions.go -prefix templates -pkg indexer definitions/

server/static.go: $(WEBSRC)
	cd web; npm run build
	go generate -v ./server

setup:
	go get -u github.com/mjibson/esc
	go get -u github.com/c4milo/github-release

install:
	go install -ldflags="$(FLAGS)" $(PREFIX)

clean-statics:
	-rm -rf web/build
	-rm server/static.go indexer/definitions.go

clean:
	-rm -rf $(BIN)*

run-dev:
	cd web/; npm start &
	rerun $(PREFIX) server --debug --passphrase "llamasrock"

docker: $(BIN)-linux-amd64
	docker-compose build

CHANNEL ?= edge
release:
	equinox release \
	--version=$(shell echo $(VERSION) | sed -e "s/^v//") \
	--config=equinox.yml \
	--channel=$(CHANNEL) \
	-- -ldflags="-X main.Version=$(VERSION) -s -w" \
	$(PREFIX)

publish:
	equinox publish \
	--release=$(shell echo $(VERSION) | sed -e "s/^v//") \
	--config=equinox.yml \
	--channel stable

github-release:
	description=$$(git cat-file -p $(VERSION) | tail -n +6); \
	commit=$$(git rev-list -n 1 $(VERSION)); \
	github-release cardigann/cardigann $(VERSION) "$$commit" "$$description" ""

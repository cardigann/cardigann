#!/bin/bash
set -eu -o pipefail

DOCKER_IMAGE=${DOCKER_IMAGE:-cardigann/cardigann}
DOCKER_TAG=${DOCKER_TAG:-$DOCKER_IMAGE:$COMMIT}
VERSION=$(git describe --tags --candidates=1)

download_cacert() {
  wget -N https://curl.haxx.se/ca/cacert.pem
}

docker_build() {
  touch server/static.go
  make cardigann-linux-amd64
  file cardigann-linux-amd64
  download_cacert
  docker build -t "${DOCKER_TAG}" .
  docker run --rm -it "${DOCKER_TAG}" version
  docker tag "${DOCKER_TAG}" "${DOCKER_IMAGE}:latest"
}

docker_login() {
  docker login -u="${DOCKER_USERNAME}" -p="${DOCKER_PASSWORD}"
}

download_cacert() {
  wget -N https://curl.haxx.se/ca/cacert.pem
}

download_equinox() {
  wget -N https://bin.equinox.io/c/mBWdkfai63v/release-tool-stable-linux-amd64.tgz
  tar -vxf release-tool-stable-linux-amd64.tgz
}

equinox_release() {
  download_equinox
  ./equinox release \
    --version="$VERSION" \
    --config ./equinox.yml \
    --channel "edge" \
    -- -ldflags="-X main.Version=$version -s -w" \
    github.com/cardigann/cardigann
}

echo "Building docker image ${DOCKER_TAG}"
docker_build
docker_login

echo "Releasing version $VERSION to equinox.io"
equinox_release "$VERSION"

echo "Pushing docker image ${DOCKER_IMAGE}"
docker push ${DOCKER_IMAGE}

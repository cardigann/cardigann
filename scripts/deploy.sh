#!/bin/bash
set -eu -o pipefail

DOCKER_IMAGE=${DOCKER_IMAGE:-cardigann/cardigann}
DOCKER_TAG=${DOCKER_TAG:-$DOCKER_IMAGE:$COMMIT}
VERSION="$(git describe --tags --candidates=1)"

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

equinox_release_edge() {
  local version="$1"
  download_equinox
  ./equinox release \
    --version="${version}" \
    --config ./equinox.yml \
    --channel edge \
    -- -ldflags="-X main.Version=${version} -s -w" \
    github.com/cardigann/cardigann
}

equinox_publish_stable() {
  local version="$1"
  ./equinox publish \
    --version="${version}" \
    --config ./equinox.yml \
    --channel stable
}

if [[ "$TRAVIS_TAG" =~ ^v ]] ; then
  download_equinox
  echo "Promoting version $TRAVIS_TAG to equinox.io stable"
  equinox_publish_stable "$TRAVIS_TAG"
  exit 0
fi

if [[ -n "$TRAVIS_TAG" ]] ; then
  echo "Skipping non-version tag"
  exit 0
fi

echo "Building docker image ${DOCKER_TAG}"
docker_build
docker_login

download_equinox

echo "Releasing version $VERSION to equinox.io edge"
equinox_release_edge "$VERSION"

echo "Pushing docker image ${DOCKER_IMAGE}"
docker push "${DOCKER_IMAGE}"

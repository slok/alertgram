#!/usr/bin/env sh

set -e

if [ -z ${VERSION} ]; then
    echo "IMAGE_VERSION env var needs to be set"
    exit 1
fi

if [ -z ${IMAGE} ]; then
    echo "IMAGE env var needs to be set"
    exit 1
fi

if [ -z ${DOCKER_FILE_PATH} ]; then
    echo "DOCKER_FILE_PATH env var needs to be set"
    exit 1
fi

docker build \
    --build-arg VERSION=${VERSION} \
    -t ${IMAGE}:${VERSION} \
    -f ${DOCKER_FILE_PATH} .
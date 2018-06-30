#!/bin/bash

VERSION=$( cat main.go | grep VERSION )
VERSION=${VERSION#*\"}
VERSION=${VERSION%\"*}
VERSION=v$VERSION

EXECUTABLE=infra-compose
DIST_DIR=dist

function d_usage() {
    echo "Usage: build.sh {build|help|release}"
}

function d_install() {
    if [[ ! -f ${GOPATH}/bin/github-release ]]; then
        echo "Install github-release"
        go get github.com/aktau/github-release
    fi
}

function d_build () {
    echo "Build crane version $VERSION"

    rm -rf $DIST_DIR
    mkdir $DIST_DIR

    GOARCH=amd64 GOOS=darwin go build -o $DIST_DIR/$EXECUTABLE-Darwin-x86_64
    GOARCH=amd64 GOOS=linux go build -o $DIST_DIR/$EXECUTABLE-Linux-x86_64
}

function d_release() {
    echo "Create release"

    d_build

    git tag $VERSION
    git push --tags

    github-release release \
        --user wayoos \
        --repo infra-compose \
        --tag $VERSION \
        --name "${VERSION}" \
        --description "infra-compose release ${VERSION}" \

    github-release upload \
        --user wayoos \
        --repo infra-compose \
        --tag $VERSION \
        --name "$EXECUTABLE-Darwin-x86_64" \
        --file $DIST_DIR/$EXECUTABLE-Darwin-x86_64

    github-release upload \
        --user wayoos \
        --repo infra-compose \
        --tag $VERSION \
        --name "$EXECUTABLE-Linux-x86_64" \
        --file $DIST_DIR/$EXECUTABLE-Linux-x86_64

}

case "$1" in
    build)
        d_build
        ;;
    install)
        d_install
        ;;
    release)
        d_release
        ;;
    help)
        d_usage
        ;;
    *)
        d_usage
        ;;
esac
#!/bin/bash

function d_usage() {
    echo "Usage: build.sh {build|help|release}"
}

function d_build () {
    goreleaser release --snapshot --skip-validate --rm-dist
}

function d_release() {
    local goreleaser_opts=""
    if [[ -z "$GITHUB_TOKEN" ]]; then
        goreleaser_opts="--skip-publish"
    fi
    goreleaser release --rm-dist $goreleaser_opts
}

case "$1" in
    build)
        d_build
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
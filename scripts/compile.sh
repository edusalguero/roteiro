#!/bin/bash
# vim: ai:ts=8:sw=8:noet
# Compile this service
# Intended to be run from local machine or CI
set -eufo pipefail
IFS=$'\t\n'

# Check required commands are in place
command -v go >/dev/null 2>&1 || { echo 'please install go or use image that has it'; exit 1; }

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd)"

if [[ -z "${TAG:-}" ]]; then
    TAG="latest"
fi

if [[ -z "${BRANCH:-}" ]]; then
    BRANCH="local-$(git rev-parse --abbrev-ref HEAD)"
fi
if [[ -z "${REPO:-}" ]] ; then
    REPO="$(basename "$(dirname "$SCRIPT_DIR")")"
fi

paths=$(find "cmd" -maxdepth 1 -mindepth 1 -name "*")

if [[ $(echo "$paths" | wc -l) -eq 0 ]]; then
  echo "No commands found"
  exit 1
fi

for path in $paths; do
    cmd=$(basename "$path")
    cmd_path="$PWD/cmd/$cmd"
    bin_path="$PWD/bin/$cmd"

    echo "Compiling $cmd"
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags netgo \
        -ldflags "\
            -X github.com/edusalguero/roteiro.git.Repo=$REPO \
            -X github.com/edusalguero/roteiro.git.Commit=$(git rev-parse --short=9 HEAD) \
            -X github.com/edusalguero/roteiro.git.Branch=$BRANCH \
            -X github.com/edusalguero/roteiro.git.Date=$(date -u '+%FT%T.%3NZ') \
            -X github.com/edusalguero/roteiro.git.Version=$TAG \
            -X \"github.com/edusalguero/roteiro.git.GoVersion=$(go version)\" \
        " \
       -o "$bin_path" "$cmd_path"
    echo "Done"
done

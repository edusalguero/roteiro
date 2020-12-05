#!/bin/bash
# vim: ai:ts=8:sw=8:noet
# Check whether the go.mod is consistent
# Intended to be run from local machine or CI
set -eufo pipefail
IFS=$'\t\n'

# Check required commands are in place
command -v go >/dev/null 2>&1 || { echo 'please install go or use image that has it'; exit 1; }

backup_go_mod()
{
    mod=$(mktemp)
    cp go.mod "$mod"

    sum=$(mktemp)
    cp go.sum "$sum"
}

restore_go_mod()
{
    cp "$mod" go.mod
    rm "$mod"

    cp "$sum" go.sum
    rm "$sum"
}

# Backup actual go.mod and go.sum
backup_go_mod
trap restore_go_mod EXIT

go mod tidy

diff "$mod" go.mod || { echo "your go.mod is inconsistent"; exit 42; }

#!/bin/bash
# vim: ai:ts=8:sw=8:noet
# Prune and add missing dependencies
# Intended to be run from local
set -eufo pipefail
IFS=$'\t\n'

# Check required commands are in place
command -v go > /dev/null 2>&1 || { echo 'please install go or use image that has it'; exit 1; }

go mod tidy
go mod verify

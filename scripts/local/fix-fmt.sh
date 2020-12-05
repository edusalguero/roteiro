#!/bin/bash
# vim: ai:ts=8:sw=8:noet
# Fix imports
# Intended to be run from local machine or CI
set -eufo pipefail
IFS=$'\t\n'

# Check required commands are in place
command -v goimports >/dev/null 2>&1 || { echo 'please install goimports or use image that has it'; exit 1; }

find . -name '*.go' -exec goimports -w {} +

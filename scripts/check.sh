#!/bin/bash
# vim: ai:ts=8:sw=8:noet
# Check this lib
# Intended to be run from local machine or CI
set -eufo pipefail
export SHELLOPTS	# propagate set to children by default
IFS=$'\t\n'

echo "lint"
bash ./scripts/checks/lint.sh
echo "test"
bash ./scripts/checks/test.sh
echo "tidy"
bash ./scripts/checks/tidy.sh
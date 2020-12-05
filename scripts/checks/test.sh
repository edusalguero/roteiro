#!/bin/bash
# vim: ai:ts=8:sw=8:noet
# Runs test suite for the service (both unit and integration tests)
set -eufo pipefail
export SHELLOPTS        # propagate set to children by default
IFS=$'\t\n'

# Check required commands are in place
command -v golangci-lint >/dev/null 2>&1 || { echo 'please install golangci-lint or use image that has it'; exit 1; }

echo "Installing dependencies"
go mod vendor

echo "-> Lint"

golangci-lint run --timeout=30m

echo "=> Lint OK"

echo "-> Running test suite"

# Setup integration test environment

go test -mod=vendor -tags=integration -race -coverprofile=.test_coverage.txt -p 1 ./...
echo "=> All tests passed"

# Collect coverage for CI report
go tool cover -func=.test_coverage.txt | tail -n1 | awk '{print "Total test coverage: " $3}'


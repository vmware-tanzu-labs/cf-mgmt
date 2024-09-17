#!/bin/bash

set -eu -o pipefail

pushd "source" >/dev/null
  go version

  echo "Running go vet"
  go vet ./...
  echo "Running go test"
  go test ./...
popd

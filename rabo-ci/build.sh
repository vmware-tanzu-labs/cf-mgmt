#!/bin/bash

set -eu -o pipefail

OUTPUT_DIR=$PWD/dist
mkdir -p ${OUTPUT_DIR}

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${OUTPUT_DIR}/cf-mgmt-linux -ldflags "-X github.com/vmwarepivotallabs/cf-mgmt/configcommands.VERSION=${VERSION} -X github.com/vmwarepivotallabs/cf-mgmt/configcommands.COMMIT=${COMMIT}" cmd/cf-mgmt/main.go
GOOS=darwin GOARCH=amd64 go build -o ${OUTPUT_DIR}/cf-mgmt-osx -ldflags "-X github.com/vmwarepivotallabs/cf-mgmt/configcommands.VERSION=${VERSION} -X github.com/vmwarepivotallabs/cf-mgmt/configcommands.COMMIT=${COMMIT}" cmd/cf-mgmt/main.go
GOOS=windows GOARCH=amd64 go build -o ${OUTPUT_DIR}/cf-mgmt.exe -ldflags "-X github.com/vmwarepivotallabs/cf-mgmt/configcommands.VERSION=${VERSION} -X github.com/vmwarepivotallabs/cf-mgmt/configcommands.COMMIT=${COMMIT}" cmd/cf-mgmt/main.go

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${OUTPUT_DIR}/cf-mgmt-config-linux -ldflags "-X github.com/vmwarepivotallabs/cf-mgmt/configcommands.VERSION=${VERSION} -X github.com/vmwarepivotallabs/cf-mgmt/configcommands.COMMIT=${COMMIT}" cmd/cf-mgmt-config/main.go
GOOS=darwin GOARCH=amd64 go build -o ${OUTPUT_DIR}/cf-mgmt-config-osx -ldflags "-X github.com/vmwarepivotallabs/cf-mgmt/configcommands.VERSION=${VERSION} -X github.com/vmwarepivotallabs/cf-mgmt/configcommands.COMMIT=${COMMIT}" cmd/cf-mgmt-config/main.go
GOOS=windows GOARCH=amd64 go build -o ${OUTPUT_DIR}/cf-mgmt-config.exe -ldflags "-X github.com/vmwarepivotallabs/cf-mgmt/configcommands.VERSION=${VERSION} -X github.com/vmwarepivotallabs/cf-mgmt/configcommands.COMMIT=${COMMIT}" cmd/cf-mgmt-config/main.go

#!/bin/bash

set -eu -o pipefail

mkdir -p ~/.ssh/ && touch ~/.ssh/known_hosts
ssh-keyscan github.com >>~/.ssh/known_hosts

SOURCE_DIR=$PWD/source

go install github.com/xchapter7x/versioning@latest

export GOPATH="$HOME/go"
export PATH="$GOPATH/bin:$PATH"

pushd ${SOURCE_DIR} > /dev/null
 if [ -d ".git" ]; then
    if ${DEV}; then
      ts=$(date +"%Y%m%M%S%N")
      DRAFT_VERSION="dev-${ts}"
      COMMIT="dev"
      VERSION="dev"
    else
      DRAFT_VERSION=`versioning bump_patch`-`git rev-parse HEAD`
      COMMIT=`git rev-parse HEAD`
      VERSION=`versioning bump_patch`
    fi
  else
    DRAFT_VERSION="v0.0.0"
    COMMIT="local"
    VERSION="v0.0.0"
  fi
  echo "next version should be: ${DRAFT_VERSION}"
popd

OUTPUT_DIR=$PWD/compiled-output
WORKING_DIR=$GOPATH/src/github.com/vmwarepivotallabs/cf-mgmt

mkdir -p ${WORKING_DIR}
cp -R ${SOURCE_DIR}/* ${WORKING_DIR}/.

pushd ${WORKING_DIR} > /dev/null
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${OUTPUT_DIR}/cf-mgmt-linux -ldflags "-X github.com/vmwarepivotallabs/cf-mgmt/configcommands.VERSION=${VERSION} -X github.com/vmwarepivotallabs/cf-mgmt/configcommands.COMMIT=${COMMIT}" cmd/cf-mgmt/main.go
  GOOS=darwin GOARCH=amd64 go build -o ${OUTPUT_DIR}/cf-mgmt-osx -ldflags "-X github.com/vmwarepivotallabs/cf-mgmt/configcommands.VERSION=${VERSION} -X github.com/vmwarepivotallabs/cf-mgmt/configcommands.COMMIT=${COMMIT}" cmd/cf-mgmt/main.go
  GOOS=windows GOARCH=amd64 go build -o ${OUTPUT_DIR}/cf-mgmt.exe -ldflags "-X github.com/vmwarepivotallabs/cf-mgmt/configcommands.VERSION=${VERSION} -X github.com/vmwarepivotallabs/cf-mgmt/configcommands.COMMIT=${COMMIT}" cmd/cf-mgmt/main.go

  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${OUTPUT_DIR}/cf-mgmt-config-linux -ldflags "-X github.com/vmwarepivotallabs/cf-mgmt/configcommands.VERSION=${VERSION} -X github.com/vmwarepivotallabs/cf-mgmt/configcommands.COMMIT=${COMMIT}" cmd/cf-mgmt-config/main.go
  GOOS=darwin GOARCH=amd64 go build -o ${OUTPUT_DIR}/cf-mgmt-config-osx -ldflags "-X github.com/vmwarepivotallabs/cf-mgmt/configcommands.VERSION=${VERSION} -X github.com/vmwarepivotallabs/cf-mgmt/configcommands.COMMIT=${COMMIT}" cmd/cf-mgmt-config/main.go
  GOOS=windows GOARCH=amd64 go build -o ${OUTPUT_DIR}/cf-mgmt-config.exe -ldflags "-X github.com/vmwarepivotallabs/cf-mgmt/configcommands.VERSION=${VERSION} -X github.com/vmwarepivotallabs/cf-mgmt/configcommands.COMMIT=${COMMIT}" cmd/cf-mgmt-config/main.go

  cp Dockerfile ${OUTPUT_DIR}/.
popd

if ${DEV}; then
  echo ${DRAFT_VERSION} > "${OUTPUT_DIR}/name"
  echo ${DRAFT_VERSION} > "${OUTPUT_DIR}/tag"
else
  echo ${VERSION} > "${OUTPUT_DIR}/name"
  echo ${VERSION} > "${OUTPUT_DIR}/tag"
fi

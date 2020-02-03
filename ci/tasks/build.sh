#!/bin/bash -e

mkdir ~/.ssh/ && touch ~/.ssh/known_hosts
ssh-keyscan github.com >>~/.ssh/known_hosts

export GOPATH=$PWD/go
export PATH=$GOPATH/bin:$PATH
OUTPUT_DIR=$PWD/compiled-output
SOURCE_DIR=$PWD/source

cp source/Dockerfile ${OUTPUT_DIR}/.

#go get github.com/Masterminds/glide
mkdir -p $GOPATH/bin
curl https://glide.sh/get | sh
go get github.com/xchapter7x/versioning

cd ${SOURCE_DIR}
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

WORKING_DIR=$GOPATH/src/github.com/pivotalservices/cf-mgmt
mkdir -p ${WORKING_DIR}
cp -R ${SOURCE_DIR}/* ${WORKING_DIR}/.
cd ${WORKING_DIR}
glide install
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${OUTPUT_DIR}/cf-mgmt-linux -ldflags "-X github.com/pivotalservices/cf-mgmt/configcommands.VERSION=${VERSION} -X github.com/pivotalservices/cf-mgmt/configcommands.COMMIT=${COMMIT}" cmd/cf-mgmt/main.go
GOOS=darwin GOARCH=amd64 go build -o ${OUTPUT_DIR}/cf-mgmt-osx -ldflags "-X github.com/pivotalservices/cf-mgmt/configcommands.VERSION=${VERSION} -X github.com/pivotalservices/cf-mgmt/configcommands.COMMIT=${COMMIT}" cmd/cf-mgmt/main.go
GOOS=windows GOARCH=amd64 go build -o ${OUTPUT_DIR}/cf-mgmt.exe -ldflags "-X github.com/pivotalservices/cf-mgmt/configcommands.VERSION=${VERSION} -X github.com/pivotalservices/cf-mgmt/configcommands.COMMIT=${COMMIT}" cmd/cf-mgmt/main.go

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${OUTPUT_DIR}/cf-mgmt-config-linux -ldflags "-X github.com/pivotalservices/cf-mgmt/configcommands.VERSION=${VERSION} -X github.com/pivotalservices/cf-mgmt/configcommands.COMMIT=${COMMIT}" cmd/cf-mgmt-config/main.go
GOOS=darwin GOARCH=amd64 go build -o ${OUTPUT_DIR}/cf-mgmt-config-osx -ldflags "-X github.com/pivotalservices/cf-mgmt/configcommands.VERSION=${VERSION} -X github.com/pivotalservices/cf-mgmt/configcommands.COMMIT=${COMMIT}" cmd/cf-mgmt-config/main.go
GOOS=windows GOARCH=amd64 go build -o ${OUTPUT_DIR}/cf-mgmt-config.exe -ldflags "-X github.com/pivotalservices/cf-mgmt/configcommands.VERSION=${VERSION} -X github.com/pivotalservices/cf-mgmt/configcommands.COMMIT=${COMMIT}" cmd/cf-mgmt-config/main.go

echo ${DRAFT_VERSION} > ${OUTPUT_DIR}/name
echo ${DRAFT_VERSION} > ${OUTPUT_DIR}/tag

#!/bin/bash -e

go get -u github.com/golang/mock/mockgen

pushd $GOPATH/src/github.com/golang/mock
  git checkout 1f837508b8ff6c01edf3bb94cc9d1fc0527f001c
  pushd mockgen
    go install
  popd
popd

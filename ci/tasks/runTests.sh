#!/bin/bash -e

mkdir ~/.ssh/ && touch ~/.ssh/known_hosts
ssh-keyscan github.com >>~/.ssh/known_hosts

export GOPATH=$PWD/go
export PATH=$GOPATH/bin:$PATH

#go get github.com/Masterminds/glide
sudo add-apt-repository ppa:masterminds/glide && sudo apt-get update -y
sudo apt-get install glide -y

WORKING_DIR=$GOPATH/src/github.com/pivotalservices/cf-mgmt
mkdir -p ${WORKING_DIR}
cp -R source/* ${WORKING_DIR}/.
cd ${WORKING_DIR}
go version
glide -v
glide install
go test $(glide nv) -v

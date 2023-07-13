#!/bin/bash -e
source ~/dev/cf-mgmt-performance/.envrc
export CONFIG_DIR=~/dev/cf-mgmt-performance/config
go run cmd/cf-mgmt/main.go export-config --excluded-org system
go run cmd/cf-mgmt/main.go apply

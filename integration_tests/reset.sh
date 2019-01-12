#!/bin/bash -e

config=config_tests/cleanup
export CONFIG_DIR=${config}
export SYSTEM_DOMAIN=dev.cfdev.sh
export USER_ID=admin
export PASSWORD=admin
export CLIENT_SECRET=admin-client-secret

config_cmd=$(mktemp)
cfmgmt_cmd=$(mktemp)
go build -o ${config_cmd} cmd/cf-mgmt-config/main.go
go build -o ${cfmgmt_cmd} cmd/cf-mgmt/main.go

${config_cmd} init
${cfmgmt_cmd} delete-orgs
rm -rf config_tests

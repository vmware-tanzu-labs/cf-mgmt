#!/bin/bash

set -eu -o pipefail

get_password_from_credhub() {
  local variable_name=$1
  credhub find -j -n "${variable_name}" | jq -r .credentials[].name | xargs credhub get -j -n | jq -r .value
}

go get code.cloudfoundry.org/uaa-cli

uaa-cli target "https://uaa.${SYSTEM_DOMAIN}"
uaa-cli get-client-credentials-token "admin" -s $(get_password_from_credhub cf_admin_password)

uaa-cli create-client cf-mgmt \
  --client_secret cf-mgmt-secret \
  --authorized_grant_types client_credentials,refresh_token \
  --authorities cloud_controller.admin,scim.read,scim.write,routing.router_groups.read

pushd source > /dev/null
  RUN_INTEGRATION_TESTS=true go test ./integration/...
popd

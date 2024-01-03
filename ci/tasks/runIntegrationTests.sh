#!/bin/bash

[ -z "$DEBUG" ] || set -x

set -eu -o pipefail

get_from_credhub() {
  local variable_name=$1
  credhub find -j -n "${variable_name}" | jq -r .credentials[].name | xargs credhub get -j -n | jq -r .value
}

get_system_domain() {
  jq -r '.cf.api_url | capture("^api\\.(?<system_domain>.*)$") | .system_domain' \
    cf-deployment-env/metadata
}

eval "$(bbl print-env --metadata-file cf-deployment-env/metadata)"

go install code.cloudfoundry.org/uaa-cli@latest

if [ -z "$SYSTEM_DOMAIN" ]
then
  SYSTEM_DOMAIN=$( get_system_domain )
fi

uaa-cli target "https://uaa.${SYSTEM_DOMAIN}" -k
uaa-cli get-client-credentials-token "admin" -s $(get_from_credhub uaa_admin_client_secret)

if ! uaa-cli get-client cf-mgmt; then
  uaa-cli create-client cf-mgmt \
    --client_secret cf-mgmt-secret \
    --authorized_grant_types client_credentials,refresh_token \
    --authorities cloud_controller.admin,scim.read,scim.write,routing.router_groups.read
fi

pushd source > /dev/null
  CF_ADMIN_PASSWORD=$(get_from_credhub cf_admin_password) \
  ADMIN_CLIENT_SECRET=$(get_from_credhub uaa_admin_client_secret) \
  RUN_INTEGRATION_TESTS=true \
    go test ./integration/... -ginkgo.progress
popd > /dev/null

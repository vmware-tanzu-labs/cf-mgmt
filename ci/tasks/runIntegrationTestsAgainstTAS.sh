#!/bin/bash

set -eu -o pipefail

[ -d env ]
[ -d source ]

: "${SYSTEM_DOMAIN:="$(jq -r '.sys_domain' env/metadata)"}"

ADMIN_CLIENT_SECRET="$( \
  om \
    --skip-ssl-validation \
    --env env/pcf.yml \
    credentials \
      --product-name cf \
      --credential-reference .uaa.admin_client_credentials \
      --credential-field password \
)"

if ! uaa-cli get-client cf-mgmt; then
  uaa-cli create-client cf-mgmt \
    --client_secret cf-mgmt-secret \
    --authorized_grant_types client_credentials,refresh_token \
    --authorities cloud_controller.admin,scim.read,scim.write,routing.router_groups.read
fi

CF_ADMIN_PASSWORD="$( \
  om \
    --skip-ssl-validation \
    --env env/pcf.yml \
    credentials \
      --product-name cf \
      --credential-reference .uaa.admin_credentials \
      --credential-field password \
)"

pushd source > /dev/null
  export ADMIN_CLIENT_SECRET CF_ADMIN_PASSWORD SYSTEM_DOMAIN

  RUN_INTEGRATION_TESTS=true \
    go test ./integration/... -ginkgo.progress
popd > /dev/null

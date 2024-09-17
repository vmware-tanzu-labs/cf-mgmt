#!/bin/bash

set -eu -o pipefail

[ -d env ]
[ -d source ]

go version

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

uaac target "uaa.$SYSTEM_DOMAIN" --skip-ssl-validation
uaac token client get admin -s "$ADMIN_CLIENT_SECRET"

uaac client delete cf-mgmt || true
uaac client add cf-mgmt \
  --name cf-mgmt \
  --secret cf-mgmt-secret \
  --authorized_grant_types client_credentials,refresh_token \
  --authorities cloud_controller.admin,scim.read,scim.write,routing.router_groups.read

CF_ADMIN_PASSWORD="$( \
  om \
    --skip-ssl-validation \
    --env env/pcf.yml \
    credentials \
    --product-name cf \
    --credential-reference .uaa.admin_credentials \
    --credential-field password \
)"

pushd source >/dev/null
  export ADMIN_CLIENT_SECRET CF_ADMIN_PASSWORD SYSTEM_DOMAIN

  RUN_INTEGRATION_TESTS=true \
    go run github.com/onsi/ginkgo/v2/ginkgo ./integration/... --show-node-events -vv --poll-progress-after
popd >/dev/null

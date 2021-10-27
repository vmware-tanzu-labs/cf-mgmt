#!/bin/bash

set -eu -o pipefail

[ -d env ]
[ -d source ]

apt-get update
apt-get install wget gnupg -y

wget -q -O - https://raw.githubusercontent.com/starkandwayne/homebrew-cf/master/public.key | apt-key add -
echo "deb http://apt.starkandwayne.com stable main" | tee /etc/apt/sources.list.d/starkandwayne.list
apt-get update
apt-get install om -y

ADMIN_CLIENT_SECRET="$( \
  om \
    --skip-ssl-validation \
    --env env/pcf.yml \
    credentials \
      --product-name cf \
      --credential-reference .uaa.admin_client_credentials \
      --credential-field password \
)"

go get code.cloudfoundry.org/uaa-cli

SYSTEM_DOMAIN="$(jq -r '.sys_domain' env/metadata)"
uaa-cli target "https://uaa.${SYSTEM_DOMAIN}" -k
uaa-cli get-client-credentials-token "admin" -s "$ADMIN_CLIENT_SECRET"

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

curl -L "https://packages.cloudfoundry.org/stable?release=linux64-binary&source=github&version=v6" | tar -zx
mv cf /usr/local/bin
cf -v

pushd source > /dev/null
  export ADMIN_CLIENT_SECRET CF_ADMIN_PASSWORD SYSTEM_DOMAIN

  RUN_INTEGRATION_TESTS=true \
    go test ./integration/... -ginkgo.progress
popd > /dev/null

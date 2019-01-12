#!/bin/bash -e

config=config_tests/many_spaces
export CONFIG_DIR=${config}
export SYSTEM_DOMAIN=dev.cfdev.sh
export USER_ID=admin
export PASSWORD=admin
export CLIENT_SECRET=admin-client-secret
export LOG_LEVEL=info

config_cmd=$(mktemp)
cfmgmt_cmd=$(mktemp)
go build -o ${config_cmd} cmd/cf-mgmt-config/main.go
go build -o ${cfmgmt_cmd} cmd/cf-mgmt/main.go

function elapsedTime() {
  start_time="$(date -u +%s)"
  $1 $2
  end_time="$(date -u +%s)"
  elapsed="$(($end_time-$start_time))"
  echo "Total of $elapsed seconds elapsed for $2"
}

${config_cmd} init
for org in {1..100}
do
  ${config_cmd} add-org --org org-${org}
  for space in {1..5}
  do
    ${config_cmd} add-space --org org-${org} --space space-${space}
  done
done

echo "Create Orgs - Run 1"
createOrg1=$(elapsedTime ${cfmgmt_cmd} create-orgs)
echo "Create Orgs - Run 2"
createOrg2=$(elapsedTime ${cfmgmt_cmd} create-orgs)
echo "Create Spaces - Run 1"
createSpace1=$(elapsedTime ${cfmgmt_cmd} create-spaces)
echo "Create Spaces - Run 2"
createSpace2=$(elapsedTime ${cfmgmt_cmd} create-spaces)
echo "Update Spaces - Run 1"
updateSpace=$(elapsedTime ${cfmgmt_cmd} update-spaces)

rm -rf config_tests
${config_cmd} init
${cfmgmt_cmd} delete-orgs
rm -rf config_tests

echo ${createOrg1}
echo ${createOrg2}
echo ${createSpace1}
echo ${createSpace2}
echo ${updateSpace}

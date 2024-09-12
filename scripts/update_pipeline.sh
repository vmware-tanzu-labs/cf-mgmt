#!/bin/bash
set -e
project_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && cd .. && pwd )"

if [[ "$(fly -t tpe-cf-mgmt status)" != "logged in successfully" ]]; then
  ${project_dir}/scripts/login_to_fly_and_save_target.sh
fi

cd ${project_dir}
fly -t tpe-cf-mgmt set-pipeline -p cf-mgmt -c <(ytt -f ci/pipelines/cf-mgmt/pipeline.yml --data-values-file ci/pipelines/cf-mgmt/values.yml)

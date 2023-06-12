#!/usr/bin/env bash
set -e
cd config-repo
cf-mgmt version
cf-mgmt $CF_MGMT_COMMAND

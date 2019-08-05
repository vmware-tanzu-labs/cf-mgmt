&larr; [back to Commands](../README.md)

# `cf-mgmt export-service-access-config`

`export-service-access-config` will capture the current visibility (all, limited, none) for all service brokers, their services and plans as well as remove the legacy configuration for this from each orgConfigs.yml

## Command Usage

```
Usage:
  cf-mgmt [OPTIONS] export-service-access-config [export-service-access-config-OPTIONS]

Help Options:
  -h, --help               Show this help message

[export-service-access-config command options]
  --config-dir=    Name of the config directory (default: config) [$CONFIG_DIR]
  --system-domain= system domain [$SYSTEM_DOMAIN]
  --user-id=       user id that has privileges to create/update/delete users, orgs and spaces [$USER_ID]
  --password=      password for user account [optional if client secret is provided] [$PASSWORD]
  --client-secret= secret for user account that has sufficient privileges to create/update/delete users, orgs and spaces] [$CLIENT_SECRET]
```

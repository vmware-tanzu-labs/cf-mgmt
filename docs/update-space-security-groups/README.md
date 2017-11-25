&larr; [back to Commands](../README.md)

# `cf-mgmt update-space-security-groups`

`update-space-security-groups` command will:
- creates/updates application security groups for a given space defined in security-group.json when `enable-security-group: true`
- assign named security groups specified in `named-security-groups: []`

## Command Usage

```
Usage:
  main [OPTIONS] update-space-security-groups [update-space-security-groups-OPTIONS]

Help Options:
  -h, --help               Show this help message

[update-space-security-groups command options]
  --config-dir=    Name of the config directory (default: config) [$CONFIG_DIR]
  --system-domain= system domain [$SYSTEM_DOMAIN]
  --user-id=       user id that has privileges to create/update/delete users, orgs and spaces [$USER_ID]
  --password=      password for user account [optional if client secret is provided] [$PASSWORD]
  --client-secret= secret for user account that has sufficient privileges to create/update/delete users,
                   orgs and spaces] [$CLIENT_SECRET]
```

&larr; [back to Commands](../README.md)

# `cf-mgmt assign-default-security-groups`

`assign-default-security-groups` command will:
- assign running security groups for anything in `running-security-groups` within cf-mgmt.yml
- assign staging security groups for anything in `staging-security-groups` within cf-mgmt.yml
- unassign any security group assigned default running or default staging that is not in `running-security-groups` or `staging-security-groups` within cf-mgmt.yml if `enable-unassign-security-groups: true`

## Command Usage
```
Usage:
  main [OPTIONS] assign-default-security-groups [assign-default-security-groups-OPTIONS]

Help Options:
  -h, --help               Show this help message

[create-security-groups command options]
  --config-dir=    Name of the config directory (default: config) [$CONFIG_DIR]
  --system-domain= system domain [$SYSTEM_DOMAIN]
  --user-id=       user id that has privileges to create/update/delete users, orgs and spaces [$USER_ID]
  --password=      password for user account [optional if client secret is provided] [$PASSWORD]
  --client-secret= secret for user account that has sufficient privileges to create/update/delete users,
                   orgs and spaces] [$CLIENT_SECRET]
  --peek           Preview entities to change without modifying. [$PEEK]
```

&larr; [back to Commands](../README.md)

# `cf-mgmt update-space-quotas`

`update-space-quotas` command will:
- creates/updates quota for a given space

## Command Usage

```
Usage:
  main [OPTIONS] update-space-quotas [update-space-quotas-OPTIONS]

Help Options:
  -h, --help               Show this help message

[update-space-quotas command options]
  --config-dir=    Name of the config directory (default: config) [$CONFIG_DIR]
  --system-domain= system domain [$SYSTEM_DOMAIN]
  --user-id=       user id that has privileges to create/update/delete users, orgs and spaces [$USER_ID]
  --password=      password for user account [optional if client secret is provided] [$PASSWORD]
  --client-secret= secret for user account that has sufficient privileges to create/update/delete users,
                   orgs and spaces] [$CLIENT_SECRET]
  --peek           Preview entities to change without modifying. [$PEEK]
```

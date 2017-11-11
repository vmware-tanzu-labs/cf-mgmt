&larr; [back to Commands](../README.md)

# `cf-mgmt update-org-quotas`

`update-org-quotas` command will:
- update org quotas specified in orgConfig.yml

## Command Usage

```
sage:
  main [OPTIONS] update-org-quotas [update-org-quotas-OPTIONS]

Help Options:
  -h, --help               Show this help message

[update-org-quotas command options]
  --config-dir=    Name of the config directory (default: config) [$CONFIG_DIR]
  --system-domain= system domain [$SYSTEM_DOMAIN]
  --user-id=       user id that has privileges to create/update/delete users, orgs and spaces [$USER_ID]
  --password=      password for user account [optional if client secret is provided] [$PASSWORD]
  --client-secret= secret for user account that has sufficient privileges to create/update/delete users,
                   orgs and spaces] [$CLIENT_SECRET]
```

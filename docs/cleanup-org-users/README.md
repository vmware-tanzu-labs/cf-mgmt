&larr; [back to Commands](../README.md)

# `cf-mgmt cleanup-org-users`

`cleanup-org-users` command will:

- will remove from any managed org any user that is not is a org or space role but was associated with the org previously

## Command Usage

```
error: Usage:
  main [OPTIONS] cleanup-org-users [cleanup-org-users-OPTIONS]

Help Options:
  -h, --help               Show this help message

[cleanup-org-users command options]
  --config-dir=    Name of the config directory (default: config) [$CONFIG_DIR]
  --system-domain= system domain [$SYSTEM_DOMAIN]
  --user-id=       user id that has privileges to create/update/delete users, orgs and spaces [$USER_ID]
  --password=      password for user account [optional if client secret is provided] [$PASSWORD]
  --client-secret= secret for user account that has sufficient privileges to create/update/delete users, orgs and spaces] [$CLIENT_SECRET]
  --peek           Preview entities to change without modifying [$PEEK]
```

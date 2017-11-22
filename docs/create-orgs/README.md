&larr; [back to Commands](../README.md)

# `cf-mgmt create-orgs`

`create-orgs` command will create all the orgs specified in orgs.yml with their configuration in orgConfig.yml

## Command Usage
```
Usage:
  main [OPTIONS] create-orgs [create-orgs-OPTIONS]

Help Options:
  -h, --help               Show this help message

[create-orgs command options]
  --config-dir=    Name of the config directory (default: config) [$CONFIG_DIR]
  --system-domain= system domain [$SYSTEM_DOMAIN]
  --user-id=       user id that has privileges to create/update/delete users, orgs and spaces [$USER_ID]
  --password=      password for user account [optional if client secret is provided] [$PASSWORD]
  --client-secret= secret for user account that has sufficient privileges to create/update/delete users,
                   orgs and spaces] [$CLIENT_SECRET]

```

&larr; [back to Commands](../README.md)

# `cf-mgmt create-org-private-domains`

`create-org-private-domains` command will:
- creates private domain(s) for a given org based on `private-domains` configured in orgConfig.yml
- Will delete any private domain(s) not in `private-domains` for a given org if `enable-remove-private-domains` is set to true

## Command Usage

```
Usage:
  main [OPTIONS] create-org-private-domains [create-org-private-domains-OPTIONS]

Help Options:
  -h, --help               Show this help message

[create-org-private-domains command options]
  --config-dir=    Name of the config directory (default: config) [$CONFIG_DIR]
  --system-domain= system domain [$SYSTEM_DOMAIN]
  --user-id=       user id that has privileges to create/update/delete users, orgs and spaces [$USER_ID]
  --password=      password for user account [optional if client secret is provided] [$PASSWORD]
  --client-secret= secret for user account that has sufficient privileges to create/update/delete users,
                   orgs and spaces] [$CLIENT_SECRET]
```

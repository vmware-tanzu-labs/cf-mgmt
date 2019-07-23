&larr; [back to Commands](../README.md)

# `cf-mgmt share-org-private-domains`

`share-org-private-domains` command will:
- shares private domain(s) for a given org based on `shared-private-domains` configured in orgConfig.yml
- Will remove (unshare) any shared private domain(s) not in `shared-private-domains` for a given org if `enable-remove-shared-private-domains` is set to true

## Command Usage

```
Usage:
  cf-mgmt [OPTIONS] share-org-private-domains [share-org-private-domains-OPTIONS]

Help Options:
  -h, --help               Show this help message

[share-org-private-domains command options]
  --config-dir=    Name of the config directory (default: config) [$CONFIG_DIR]
  --system-domain= system domain [$SYSTEM_DOMAIN]
  --user-id=       user id that has privileges to create/update/delete users, orgs and spaces [$USER_ID]
  --password=      password for user account [optional if client secret is provided] [$PASSWORD]
  --client-secret= secret for user account that has sufficient privileges to create/update/delete users,
                   orgs and spaces] [$CLIENT_SECRET]
  --peek           Preview entities to change without modifying. [$PEEK]
```

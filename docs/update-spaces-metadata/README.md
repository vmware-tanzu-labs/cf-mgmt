&larr; [back to Commands](../README.md)

# `cf-mgmt update-spaces-metadata`

`update-spaces-metadata` command will:

- add metadata label for given space
- add metadata annotation for given space
- remove metadata label for given space
- remove metadata annotation for given space
- remove metadata from given space

## Command Usage

```
Usage:
  cf-mgmt [OPTIONS] update-spaces-metadata [update-spaces-metadata-OPTIONS]

Help Options:
  -h, --help               Show this help message

[update-spaces-metadata command options]
  --config-dir=    Name of the config directory (default: config) [$CONFIG_DIR]
  --system-domain= system domain [$SYSTEM_DOMAIN]
  --user-id=       user id that has privileges to create/update/delete users, orgs and spaces [$USER_ID]
  --password=      password for user account [optional if client secret is provided] [$PASSWORD]
  --client-secret= secret for user account that has sufficient privileges to create/update/delete users,
                   orgs and spaces] [$CLIENT_SECRET]
  --ldap-password= LDAP password for binding [$LDAP_PASSWORD]
  --peek           Preview entities to change without modifying. [$PEEK]
```

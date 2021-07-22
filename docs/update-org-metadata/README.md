&larr; [back to Commands](../README.md)

# `cf-mgmt update-orgs-metadata`

Note: if no `metadata-prefix` is provided in `cf-mgmt.yml` this command will
default to using `cf-mgmt.pivotal.io` as the prefix.

`update-org-metadata` command will:

- add metadata label for given org
- add metadata annotation for given org
- remove metadata label for given org
- remove metadata annotation for given org
- remove metadata from given org

## Command Usage

```
Usage:
  cf-mgmt [OPTIONS] update-orgs-metadata [update-orgs-metadata-OPTIONS]

Help Options:
  -h, --help               Show this help message

[update-orgs-metadata command options]
  --config-dir=    Name of the config directory (default: config) [$CONFIG_DIR]
  --system-domain= system domain [$SYSTEM_DOMAIN]
  --user-id=       user id that has privileges to create/update/delete users, orgs and spaces [$USER_ID]
  --password=      password for user account [optional if client secret is provided] [$PASSWORD]
  --client-secret= secret for user account that has sufficient privileges to create/update/delete users,
                   orgs and spaces] [$CLIENT_SECRET]
  --ldap-password= LDAP password for binding [$LDAP_PASSWORD]
  --peek           Preview entities to change without modifying. [$PEEK]
```

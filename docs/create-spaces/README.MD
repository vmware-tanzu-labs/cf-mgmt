&larr; [back to Commands](../README.md)

# `cf-mgmt create-spaces`

`create-spaces` command will:
- creates spaces for all spaces listed in each spaces.yml

## Command Usage

```
Usage:
  main [OPTIONS] create-spaces [create-spaces-OPTIONS]

Help Options:
  -h, --help               Show this help message

[create-spaces command options]
  --config-dir=    Name of the config directory (default: config) [$CONFIG_DIR]
  --system-domain= system domain [$SYSTEM_DOMAIN]
  --user-id=       user id that has privileges to create/update/delete users, orgs and spaces [$USER_ID]
  --password=      password for user account [optional if client secret is provided] [$PASSWORD]
  --client-secret= secret for user account that has sufficient privileges to create/update/delete users,
                   orgs and spaces] [$CLIENT_SECRET]
  --ldap-password= LDAP password for binding [$LDAP_PASSWORD]
```

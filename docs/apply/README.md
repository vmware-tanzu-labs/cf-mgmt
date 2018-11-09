&larr; [back to Commands](../README.md)

# `cf-mgmt apply`

`apply` will run all the commands in the correct order.  Ideal for using peek option to see what will happen or if not using concourse pipeline.

## Command Usage

```
error: Usage:
  main [OPTIONS] apply [apply-OPTIONS]

Help Options:
  -h, --help               Show this help message

[apply command options]
  --config-dir=    Name of the config directory (default: config) [$CONFIG_DIR]
  --system-domain= system domain [$SYSTEM_DOMAIN]
  --user-id=       user id that has privileges to create/update/delete users, orgs and spaces [$USER_ID]
  --password=      password for user account [optional if client secret is provided] [$PASSWORD]
  --client-secret= secret for user account that has sufficient privileges to create/update/delete users, orgs and spaces] [$CLIENT_SECRET]
  --peek           Preview entities to change without modifying [$PEEK]
  --ldap-password= LDAP password for binding [$LDAP_PASSWORD]
```

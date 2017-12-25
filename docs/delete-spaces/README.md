&larr; [back to Commands](../README.md)

# `cf-mgmt delete-spaces`

`delete-spaces` command will:
- if `enable-delete-spaces: true` in spaces.yml it will delete spaces NOT specified in spaces.yml.  This is recursive delete apps, brokers, etc for a space.
- specifying `--peek` will show you which spaces would be deleted, without actually deleting them.

## Command Usage

```
Usage:
  main [OPTIONS] delete-spaces [delete-spaces-OPTIONS]

Help Options:
  -h, --help               Show this help message

[delete-spaces command options]
  --config-dir=    Name of the config directory (default: config) [$CONFIG_DIR]
  --system-domain= system domain [$SYSTEM_DOMAIN]
  --user-id=       user id that has privileges to create/update/delete users, orgs and spaces [$USER_ID]
  --password=      password for user account [optional if client secret is provided] [$PASSWORD]
  --client-secret= secret for user account that has sufficient privileges to create/update/delete users,
                   orgs and spaces] [$CLIENT_SECRET]
  --peek           Preview entities to change without modifying. [$PEEK]
```

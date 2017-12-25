&larr; [back to Commands](../README.md)

# `cf-mgmt delete-orgs`

`delete-orgs` command will delete orgs from your cloud foundry installation
- deletes orgs NOT specified in orgs.yml.  This is recursive for underlaying spaces and apps.
- Will NOT delete orgs which are `protected_orgs` in orgs.yml
- specifying `--peek` will show you which orgs would be deleted, without actually deleting them.

## Command Usage

```
Usage:
  main [OPTIONS] delete-orgs [delete-orgs-OPTIONS]

Help Options:
  -h, --help               Show this help message

[delete-orgs command options]
  --config-dir=    Name of the config directory (default: config) [$CONFIG_DIR]
  --system-domain= system domain [$SYSTEM_DOMAIN]
  --user-id=       user id that has privileges to create/update/delete users, orgs and spaces [$USER_ID]
  --password=      password for user account [optional if client secret is provided] [$PASSWORD]
  --client-secret= secret for user account that has sufficient privileges to create/update/delete users,
                   orgs and spaces] [$CLIENT_SECRET]
  --peek           Preview entities to change without modifying. [$PEEK]
```

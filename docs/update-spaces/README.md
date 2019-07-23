&larr; [back to Commands](../README.md)

# `cf-mgmt update-spaces`

`update-spaces` command will:

- updates allow ssh property true/false for given spaceConfig.yml

## Command Usage

```
Usage:
  cf-mgmt [OPTIONS] update-spaces [update-spaces-OPTIONS]

Help Options:
  -h, --help               Show this help message

[update-spaces command options]
  --config-dir=    Name of the config directory (default: config) [$CONFIG_DIR]
  --system-domain= system domain [$SYSTEM_DOMAIN]
  --user-id=       user id that has privileges to create/update/delete users, orgs and spaces [$USER_ID]
  --password=      password for user account [optional if client secret is provided] [$PASSWORD]
  --client-secret= secret for user account that has sufficient privileges to create/update/delete users,
                   orgs and spaces] [$CLIENT_SECRET]
  --peek           Preview entities to change without modifying. [$PEEK]
```

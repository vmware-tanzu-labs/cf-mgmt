&larr; [back to Commands](../README.md)

# `cf-mgmt service-access`

`service-access` command will:
- disable public access on all services
- enable access to all plans for any org that is in `protected-orgs`
- will enable access for a given service and plan(s) based on configuration

Included with v1.0.7+

## Command Usage

```
Usage:
  cf-mgmt [OPTIONS] service-access [service-access-OPTIONS]

Help Options:
  -h, --help               Show this help message

[service-access command options]
  --config-dir=    Name of the config directory (default: config) [$CONFIG_DIR]
  --system-domain= system domain [$SYSTEM_DOMAIN]
  --user-id=       user id that has privileges to create/update/delete users, orgs and spaces [$USER_ID]
  --password=      password for user account [optional if client secret is provided] [$PASSWORD]
  --client-secret= secret for user account that has sufficient privileges to create/update/delete users, orgs and spaces] [$CLIENT_SECRET]
  --peek           Preview entities to change without modifying [$PEEK]
```

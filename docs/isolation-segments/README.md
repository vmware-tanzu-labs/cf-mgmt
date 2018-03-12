&larr; [back to Commands](../README.md)

# `cf-mgmt isolation-segments`

`isolation-segments` command will:
- ensure that all isolation segments that are defined in `default_isolation_segment` for any orgConfig.yml is defined
- ensure that any spaces is associated with isolation segments in `isolation_segment` field of spaceConfig.yml
- will remove isolation segment definitions not in configuration if `enable-delete-isolation-segments: true` in cf-mgmt.yml in the config directory

**Note: isolation segment must be deployed by platform engineering team matching names used***

## Command Usage
```
Usage:
  main [OPTIONS] isolation-segments [isolation-segments-OPTIONS]

Help Options:
  -h, --help               Show this help message

[isolation-segments command options]
  --config-dir=    Name of the config directory (default: config) [$CONFIG_DIR]
  --system-domain= system domain [$SYSTEM_DOMAIN]
  --user-id=       user id that has privileges to create/update/delete users, orgs and spaces [$USER_ID]
  --password=      password for user account [optional if client secret is provided] [$PASSWORD]
  --client-secret= secret for user account that has sufficient privileges to create/update/delete users, orgs and spaces] [$CLIENT_SECRET]
  --peek           Preview entities to change without modifying [$PEEK]
```

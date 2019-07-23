&larr; [back to Commands](../README.md)

# `cf-mgmt create-security-groups`

`create-security-groups` command will:
- create a named asg for any .json file in asgs folder in root of config directory.  These asgs will be named based on file name mysql.json will create an asg named mysql.
- create a named asg for any .json file in default_asgs folder in root of config directory.  These asgs are meant to be for running and staging default asgs

## Command Usage
```
Usage:
  cf-mgmt [OPTIONS] create-security-groups [create-security-groups-OPTIONS]

Help Options:
  -h, --help               Show this help message

[create-security-groups command options]
  --config-dir=    Name of the config directory (default: config) [$CONFIG_DIR]
  --system-domain= system domain [$SYSTEM_DOMAIN]
  --user-id=       user id that has privileges to create/update/delete users, orgs and spaces [$USER_ID]
  --password=      password for user account [optional if client secret is provided] [$PASSWORD]
  --client-secret= secret for user account that has sufficient privileges to create/update/delete users,
                   orgs and spaces] [$CLIENT_SECRET]
  --peek           Preview entities to change without modifying. [$PEEK]
```

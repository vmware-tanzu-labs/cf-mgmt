&larr; [back to Commands](../README.md)

# `cf-mgmt export-config`

`export-config` will export org/space/user details from an existing Cloud Foundry instance. This is useful when you have an existing foundation and would like to use the `cf-mgmt` to manage that foundation.

Once your run `./cf-mgmt export-config`, a config directory with org and space details will be created. This will also export user details such as org and space users and their roles within specific org and space. Other details exported include org and space quota details and ssh access at space level.

You can exclude orgs and spaces from export by using the flag `--excluded-org` and for space `--excluded-space`.

```
WARNING : Running this command will delete existing config folder and will create it again with the new configuration
```

`NOTE: Please make sure to enable and configure LDAP after export if your foundation is ldap enabled. Otherwise when the pipeline runs, it will un map the user roles assuming that they don't exists in LDAP`

## Command Usage

```
Usage:
  cf-mgmt [OPTIONS] export-config [export-config-OPTIONS]

Help Options:
  -h, --help                Show this help message

[export-config command options]
  --config-dir=     Name of the config directory (default: config) [$CONFIG_DIR]
  --system-domain=  system domain [$SYSTEM_DOMAIN]
  --user-id=        user id that has privileges to create/update/delete users, orgs and spaces [$USER_ID]
  --password=       password for user account [optional if client secret is provided] [$PASSWORD]
  --client-secret=  secret for user account that has sufficient privileges to create/update/delete users,
                    orgs and spaces] [$CLIENT_SECRET]
  --excluded-org=   Org to be excluded from export. Repeat the flag to specify multiple orgs
  --excluded-space= Space to be excluded from export. Repeat the flag to specify multiple spaces
```

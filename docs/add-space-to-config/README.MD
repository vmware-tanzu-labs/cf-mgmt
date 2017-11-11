&larr; [back to Commands](../README.md)

# `cf-mgmt add-space-to-config`

`add-space-to-config` allows for adding a space to a previously defined org.  This will generate a folder for each space inside the orgs folder.  In the spaces folder will contain a spaceConfig.yml and a security-group.json file.  Any space listed in spaces.yml will be created when the create-spaces operation is ran.  

## Command Usage

```
Usage:
  main [OPTIONS] add-space-to-config [add-space-to-config-OPTIONS]

Help Options:
  -h, --help                   Show this help message

[add-space-to-config command options]
  --config-dir=        Name of the config directory (default: config) [$CONFIG_DIR]
  --org=               Org name to add [$ORG]
  --space=             Space name to add [$space]
  --space-dev-grp=     LDAP group for Space Developer [$SPACE_DEV_GRP]
  --space-mgr-grp=     LDAP group for Space Manager [$SPACE_MGR_GRP]
  --space-auditor-grp= LDAP group for Space Auditor [$SPACE_AUDITOR_GRP]
```

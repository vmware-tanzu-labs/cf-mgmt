&larr; [back to Commands](../README.md)

# `cf-mgmt add-org-to-config`

`add-org-to-config` will add the specified org to orgs.yml and create folder based on the org name you specified.  Within this folder will be an orgConfig.yml and spaces.yml which will be empty but will eventually contain a list of spaces.  Any org listed in orgs.yml will be created when the create-orgs operation is ran.

## Command Usage

```
Usage:
  main [OPTIONS] add-org-to-config [add-org-to-config-OPTIONS]

Help Options:
  -h, --help                     Show this help message

[add-org-to-config command options]
  --config-dir=          Name of the config directory (default: config) [$CONFIG_DIR]
  --org=                 Org name to add [$ORG]
  --org-billing-mgr-grp= LDAP group for Org Billing Manager [$ORG_BILLING_MGR_GRP]
  --org-mgr-grp=         LDAP group for Org Manager [$ORG_MGR_GRP]
  --org-auditor-grp=     LDAP group for Org Auditor [$ORG_AUDITOR_GRP]
```

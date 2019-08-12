&larr; [back to Commands](../README.md)

# `cf-mgmt-config add-org`

*** Deprecated *** - Use `org` command instead

`add-org` will add the specified org to orgs.yml and create folder based on the org name you specified.  Within this folder will be an orgConfig.yml and spaces.yml which will be empty but will eventually contain a list of spaces.  Any org listed in orgs.yml will be created when the create-orgs operation is ran.

## Command Usage

```
Usage:
  cf-mgmt-config [OPTIONS] add-org [add-org-OPTIONS]

Help Options:
  -h, --help                                        Show this help message

[add-org command options]
  --config-dir=                             Name of the config directory (default: config) [$CONFIG_DIR]
  --org=                                    Org name
  --private-domain=                         Private Domain(s) to add, specify multiple times
  --shared-private-domain=                  Shared Private Domain(s) to add, specify multiple times
  --default-isolation-segment=              Default isolation segment for org
  --named-quota=                            Named quota to assign to org
  --enable-remove-spaces=[true|false]       Enable removing spaces

quota:
  --enable-org-quota=[true|false]           Enable the Org Quota in the config
  --memory-limit=                           An Org's memory limit in Megabytes
  --instance-memory-limit=                  Global Org Application instance memory limit in Megabytes
  --total-routes=                           Total Routes capacity for an Org
  --total-services=                         Total Services capacity for an Org
  --paid-service-plans-allowed=[true|false] Allow paid services to appear in an org
  --total-private-domains=                  Total Private Domain capacity for an Org
  --total-reserved-route-ports=             Total Reserved Route Ports capacity for an Org
  --total-service-keys=                     Total Service Keys capacity for an Org
  --app-instance-limit=                     Total Service Keys capacity for an Org

billing-manager:
  --billing-manager-ldap-user=              Ldap User to add, specify multiple times
  --billing-manager-user=                   User to add, specify multiple times
  --billing-manager-saml-user=              SAML user to add, specify multiple times
  --billing-manager-ldap-group=             Group to add, specify multiple times

manager:
  --manager-ldap-user=                      Ldap User to add, specify multiple times
  --manager-user=                           User to add, specify multiple times
  --manager-saml-user=                      SAML user to add, specify multiple times
  --manager-ldap-group=                     Group to add, specify multiple times

auditor:
  --auditor-ldap-user=                      Ldap User to add, specify multiple times
  --auditor-user=                           User to add, specify multiple times
  --auditor-saml-user=                      SAML user to add, specify multiple times
  --auditor-ldap-group=                     Group to add, specify multiple times
service-access:
  --service=                                *** Deprecated *** Service Name to add, specify multiple times
```

&larr; [back to Commands](../README.md)

# `cf-mgmt-config update-org`

`update-org` command allows updating any property except name of org within orgConfig.yml
- quotas
- user/role mappings

## Command Usage
```
Usage:
  main [OPTIONS] update-org [update-org-OPTIONS]

Help Options:
  -h, --help                                           Show this help message

[update-org command options]
  --config-dir=                                       Name of the config directory (default: config) [$CONFIG_DIR]
  --org=                                              Org name
  --private-domain=                                   Private Domain(s) to add, specify multiple times
  --private-domain-to-remove=                         Private Domain(s) to remove, specify multiple times
  --enable-remove-private-domains=[true|false]        Enable removing private domains
  --shared-private-domain=                            Shared Private Domain(s) to add, specify multiple times
  --shared-private-domain-to-remove=                  Shared Private Domain(s) to remove, specify multiple times
  --enable-remove-shared-private-domains=[true|false] Enable removing shared private domains
  --default-isolation-segment=                        Default isolation segment for org
  --clear-default-isolation-segment                   Sets the default isolation segment to blank
  --enable-remove-users=[true|false]                  Enable removing users from the org

quota:
  --enable-org-quota=[true|false]              Enable the Org Quota in the config
  --memory-limit=                              An Org's memory limit in Megabytes
  --instance-memory-limit=                     Global Org Application instance memory limit in Megabytes
  --total-routes=                              Total Routes capacity for an Org
  --total-services=                            Total Services capacity for an Org
  --paid-service-plans-allowed=[true|false]    Allow paid services to appear in an org
  --total-private-domains=                     Total Private Domain capacity for an Org
  --total-reserved-route-ports=                Total Reserved Route Ports capacity for an Org
  --total-service-keys=                        Total Service Keys capacity for an Org
  --app-instance-limit=                        Total Service Keys capacity for an Org

billing-manager:
  --billing-manager-ldap-user=                 Ldap User to add, specify multiple times
  --billing-manager-ldap-user-to-remove=       Ldap User to remove, specify multiple times
  --billing-manager-user=                      User to add, specify multiple times
  --billing-manager-user-to-remove=            User to remove, specify multiple times
  --billing-manager-saml-user=                 SAML user to add, specify multiple times
  --billing-manager-saml-user-to-remove=       SAML user to remove, specify multiple times
  --billing-manager-ldap-group=                Group to add, specify multiple times
  --billing-manager-ldap-group-to-remove=      Group to remove, specify multiple times

manager:
  --manager-ldap-user=                         Ldap User to add, specify multiple times
  --manager-ldap-user-to-remove=               Ldap User to remove, specify multiple times
  --manager-user=                              User to add, specify multiple times
  --manager-user-to-remove=                    User to remove, specify multiple times
  --manager-saml-user=                         SAML user to add, specify multiple times
  --manager-saml-user-to-remove=               SAML user to remove, specify multiple times
  --manager-ldap-group=                        Group to add, specify multiple times
  --manager-ldap-group-to-remove=              Group to remove, specify multiple times

auditor:
  --auditor-ldap-user=                         Ldap User to add, specify multiple times
  --auditor-ldap-user-to-remove=               Ldap User to remove, specify multiple times
  --auditor-user=                              User to add, specify multiple times
  --auditor-user-to-remove=                    User to remove, specify multiple times
  --auditor-saml-user=                         SAML user to add, specify multiple times
  --auditor-saml-user-to-remove=               SAML user to remove, specify multiple times
  --auditor-ldap-group=                        Group to add, specify multiple times
  --auditor-ldap-group-to-remove=              Group to remove, specify multiple times
```

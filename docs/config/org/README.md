&larr; [back to Commands](../README.md)

# `cf-mgmt-config org`

`org` command allows adding/updating any property except name of org name within orgConfig.yml

- quotas
- user/role mappings
- service access

## Command Usage

```
Usage:
  cf-mgmt-config [OPTIONS] org [org-OPTIONS]

Help Options:
  -h, --help                                                  Show this help message

[org command options]
      --config-dir=                                       Name of the config directory (default: config) [$CONFIG_DIR]
      --org=                                              Org name
      --private-domain=                                   Private Domain(s) to add, specify multiple times
      --private-domain-to-remove=                         Private Domain(s) to remove, specify multiple times
      --enable-remove-private-domains=[true|false]        Enable removing private domains
      --shared-private-domain=                            Shared Private Domain(s) to add, specify multiple times
      --shared-private-domain-to-remove=                  Shared Private Domain(s) to remove, specify multiple times
      --enable-remove-shared-private-domains=[true|false] Enable removing shared private domains
      --enable-remove-spaces=[true|false]                 Enable removing spaces
      --default-isolation-segment=                        Default isolation segment for org
      --clear-default-isolation-segment                   Sets the default isolation segment to blank
      --enable-remove-users=[true|false]                  Enable removing users from the org
      --named-quota=                                      Named quota to assign to org
      --clear-named-quota                                 Sets the named quota to blank

quota:
      --enable-org-quota=[true|false]                     Enable the Org Quota in the config
      --memory-limit=                                     An Org's memory limit in Megabytes
      --instance-memory-limit=                            Global Org Application instance memory limit in Megabytes
      --total-routes=                                     Total Routes capacity for an Org
      --total-services=                                   Total Services capacity for an Org
      --paid-service-plans-allowed=[true|false]           Allow paid services to appear in an org
      --total-private-domains=                            Total Private Domain capacity for an Org
      --total-reserved-route-ports=                       Total Reserved Route Ports capacity for an Org
      --total-service-keys=                               Total Service Keys capacity for an Org
      --app-instance-limit=                               App Instance Limit an Org
      --app-task-limit=                                   App Task Limit an Org
      --log-rate-limit-bytes-per-second=                  Log Rate limit per app for an org

billing-manager:
      --billing-manager-ldap-user=                        Ldap User to add, specify multiple times
      --billing-manager-user=                             User to add, specify multiple times
      --billing-manager-saml-user=                        SAML user to add, specify multiple times
      --billing-manager-ldap-group=                       Group to add, specify multiple times
      --billing-manager-ldap-user-to-remove=              Ldap User to remove, specify multiple times
      --billing-manager-user-to-remove=                   User to remove, specify multiple times
      --billing-manager-saml-user-to-remove=              SAML user to remove, specify multiple times
      --billing-manager-ldap-group-to-remove=             Group to remove, specify multiple times

manager:
      --manager-ldap-user=                                Ldap User to add, specify multiple times
      --manager-user=                                     User to add, specify multiple times
      --manager-saml-user=                                SAML user to add, specify multiple times
      --manager-ldap-group=                               Group to add, specify multiple times
      --manager-ldap-user-to-remove=                      Ldap User to remove, specify multiple times
      --manager-user-to-remove=                           User to remove, specify multiple times
      --manager-saml-user-to-remove=                      SAML user to remove, specify multiple times
      --manager-ldap-group-to-remove=                     Group to remove, specify multiple times

auditor:
      --auditor-ldap-user=                                Ldap User to add, specify multiple times
      --auditor-user=                                     User to add, specify multiple times
      --auditor-saml-user=                                SAML user to add, specify multiple times
      --auditor-ldap-group=                               Group to add, specify multiple times
      --auditor-ldap-user-to-remove=                      Ldap User to remove, specify multiple times
      --auditor-user-to-remove=                           User to remove, specify multiple times
      --auditor-saml-user-to-remove=                      SAML user to remove, specify multiple times
      --auditor-ldap-group-to-remove=                     Group to remove, specify multiple times

service-access:
      --service=                                          *****DEPRECATED, use 'cf-mgmt-config global service-access' ***** - Service Name to add
      --plans=                                            *****DEPRECATED, use 'cf-mgmt-config global service-access' ***** - plans to add, empty list will add all plans
      --service-to-remove=                                *****DEPRECATED, use 'cf-mgmt-config global service-access' ***** - name of service to remove

metadata:
      --label=                                      Label to add, can specify multiple
      --label-value=                                Label value to add, can specify multiple but need to match number of label args
      --annotation=                                 Annotation to add, can specify multiple
      --annotation-value=                           Annotation value to add, can specify multiple but need to match number of annotation args
      --labels-to-remove=                           name of label to remove
      --annotations-to-remove=                      name of annotation to remove
```

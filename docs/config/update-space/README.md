&larr; [back to Commands](../README.md)

# `cf-mgmt-config update-space`

`update-space` command allows updating any property except name of org/space within spaceConfig.yml
- quotas
- user/role mappings

## Command Usage
```
Usage:
  main [OPTIONS] update-space [update-space-OPTIONS]

Help Options:
  -h, --help                                        Show this help message

[update-space command options]
  --config-dir=                             Name of the config directory (default: config) [$CONFIG_DIR]
  --org=                                    Org name
  --space=                                  Space name
  --allow-ssh=[true|false]                  Enable the Space Quota in the config
  --enable-remove-users=[true|false]        Enable removing users from the space
  --isolation-segment=                      Isolation segment assigned to space
  --clear-isolation-segment                 Sets the isolation segment to blank
  --named-asg=                              Named asg(s) to assign to space, specify multiple times
  --named-asg-to-remove=                    Named asg(s) to remove, specify multiple times

quota:
  --enable-space-quota=[true|false]         Enable the Space Quota in the config
  --memory-limit=                           An Space's memory limit in Megabytes
  --instance-memory-limit=                  Space Application instance memory limit in Megabytes
  --total-routes=                           Total Routes capacity for an Space
  --total-services=                         Total Services capacity for an Space
  --paid-service-plans-allowed=[true|false] Allow paid services to appear in an Space
  --total-private-domains=                  Total Private Domain capacity for an Space
  --total-reserved-route-ports=             Total Reserved Route Ports capacity for an Space
  --total-service-keys=                     Total Service Keys capacity for an Space
  --app-instance-limit=                     Total Service Keys capacity for an Space

developer:
  --developer-ldap-user=                    Ldap User to add, specify multiple times
  --developer-ldap-user-to-remove=          Ldap User to remove, specify multiple times
  --developer-user=                         User to add, specify multiple times
  --developer-user-to-remove=               User to remove, specify multiple times
  --developer-saml-user=                    SAML user to add, specify multiple times
  --developer-saml-user-to-remove=          SAML user to remove, specify multiple times
  --developer-ldap-group=                   Group to add, specify multiple times
  --developer-ldap-group-to-remove=         Group to remove, specify multiple times

manager:
  --manager-ldap-user=                      Ldap User to add, specify multiple times
  --manager-ldap-user-to-remove=            Ldap User to remove, specify multiple times
  --manager-user=                           User to add, specify multiple times
  --manager-user-to-remove=                 User to remove, specify multiple times
  --manager-saml-user=                      SAML user to add, specify multiple times
  --manager-saml-user-to-remove=            SAML user to remove, specify multiple times
  --manager-ldap-group=                     Group to add, specify multiple times
  --manager-ldap-group-to-remove=           Group to remove, specify multiple times

auditor:
  --auditor-ldap-user=                      Ldap User to add, specify multiple times
  --auditor-ldap-user-to-remove=            Ldap User to remove, specify multiple times
  --auditor-user=                           User to add, specify multiple times
  --auditor-user-to-remove=                 User to remove, specify multiple times
  --auditor-saml-user=                      SAML user to add, specify multiple times
  --auditor-saml-user-to-remove=            SAML user to remove, specify multiple times
  --auditor-ldap-group=                     Group to add, specify multiple times
  --auditor-ldap-group-to-remove=           Group to remove, specify multiple times
```

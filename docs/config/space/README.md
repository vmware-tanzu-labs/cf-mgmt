&larr; [back to Commands](../README.md)

# `cf-mgmt-config space`

`space` command allows adding/updating any property except name of org/space within spaceConfig.yml

- quotas
- user/role mappings

## Command Usage

```
Usage:
  cf-mgmt-config [OPTIONS] space [space-OPTIONS]

Help Options:
  -h, --help                                            Show this help message

[space command options]
      --config-dir=                                 Name of the config directory (default: config) [$CONFIG_DIR]
      --org=                                        Org name
      --space=                                      Space name
      --allow-ssh=[true|false]                      Enable the application ssh
      --allow-ssh-until=                            Temporarily allow application ssh until options are Days (1D), Hours (5H), or Minutes (10M)
      --enable-remove-users=[true|false]            Enable removing users from the space
      --enable-security-group=[true|false]          Enable space level security group definitions
      --enable-unassign-security-group=[true|false] Enable unassigning security groups not in config
      --isolation-segment=                          Isolation segment assigned to space
      --clear-isolation-segment                     Sets the isolation segment to blank
      --named-asg=                                  Named asg(s) to assign to space, specify multiple times
      --named-asg-to-remove=                        Named asg(s) to remove, specify multiple times
      --named-quota=                                Named quota to assign to space
      --clear-named-quota                           Sets the named quota to blank

quota:
      --enable-space-quota=[true|false]             Enable the Space Quota in the config
      --memory-limit=                               An Space's memory limit in Megabytes
      --instance-memory-limit=                      Space Application instance memory limit in Megabytes
      --total-routes=                               Total Routes capacity for an Space
      --total-services=                             Total Services capacity for an Space
      --paid-service-plans-allowed=[true|false]     Allow paid services to appear in an Space
      --total-reserved-route-ports=                 Total Reserved Route Ports capacity for an Space
      --total-service-keys=                         Total Service Keys capacity for an Space
      --app-instance-limit=                         App Instance Limit for a space
      --app-task-limit=                             App Task Limit for a space
      --log-rate-limit-bytes-per-second=            Log Rate limit per app for a space

developer:
      --developer-ldap-user=                        Ldap User to add, specify multiple times
      --developer-user=                             User to add, specify multiple times
      --developer-saml-user=                        SAML user to add, specify multiple times
      --developer-ldap-group=                       Group to add, specify multiple times
      --developer-ldap-user-to-remove=              Ldap User to remove, specify multiple times
      --developer-user-to-remove=                   User to remove, specify multiple times
      --developer-saml-user-to-remove=              SAML user to remove, specify multiple times
      --developer-ldap-group-to-remove=             Group to remove, specify multiple times

manager:
      --manager-ldap-user=                          Ldap User to add, specify multiple times
      --manager-user=                               User to add, specify multiple times
      --manager-saml-user=                          SAML user to add, specify multiple times
      --manager-ldap-group=                         Group to add, specify multiple times
      --manager-ldap-user-to-remove=                Ldap User to remove, specify multiple times
      --manager-user-to-remove=                     User to remove, specify multiple times
      --manager-saml-user-to-remove=                SAML user to remove, specify multiple times
      --manager-ldap-group-to-remove=               Group to remove, specify multiple times

auditor:
      --auditor-ldap-user=                          Ldap User to add, specify multiple times
      --auditor-user=                               User to add, specify multiple times
      --auditor-saml-user=                          SAML user to add, specify multiple times
      --auditor-ldap-group=                         Group to add, specify multiple times
      --auditor-ldap-user-to-remove=                Ldap User to remove, specify multiple times
      --auditor-user-to-remove=                     User to remove, specify multiple times
      --auditor-saml-user-to-remove=                SAML user to remove, specify multiple times
      --auditor-ldap-group-to-remove=               Group to remove, specify multiple times

supporter:
      --supporter-ldap-user=                          Ldap User to add, specify multiple times
      --supporter-user=                               User to add, specify multiple times
      --supporter-saml-user=                          SAML user to add, specify multiple times
      --supporter-ldap-group=                         Group to add, specify multiple times
      --supporter-ldap-user-to-remove=                Ldap User to remove, specify multiple times
      --supporter-user-to-remove=                     User to remove, specify multiple times
      --supporter-saml-user-to-remove=                SAML user to remove, specify multiple times
      --supporter-ldap-group-to-remove=               Group to remove, specify multiple times

metadata:
      --label=                                      Label to add, can specify multiple
      --label-value=                                Label value to add, can specify multiple but need to match number of label args
      --annotation=                                 Annotation to add, can specify multiple
      --annotation-value=                           Annotation value to add, can specify multiple but need to match number of annotation args
      --labels-to-remove=                           name of label to remove
      --annotations-to-remove=                      name of annotation to remove
```

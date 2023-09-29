&larr; [back to Commands](../README.md)

# `cf-mgmt-config add-space`

**_ Deprecated _** - Use `space` command instead

`add-space` allows for adding a space to a previously defined org. This will generate a folder for each space inside the orgs folder. In the spaces folder will contain a spaceConfig.yml and a security-group.json file. Any space listed in spaces.yml will be created when the create-spaces operation is ran.

## Command Usage

```
Usage:
  cf-mgmt-config [OPTIONS] add-space [add-space-OPTIONS]

Help Options:
  -h, --help                                        Show this help message

[add-space command options]
  --config-dir=                             Name of the config directory (default: config) [$CONFIG_DIR]
  --org=                                    Org name
  --space=                                  Space name
  --allow-ssh=[true|false]                  Enable the application ssh
  --allow-ssh-until=                        Temporarily allow application ssh until options are Days (1D), Hours (5H), or Minutes (10M)
  --enable-security-group=[true|false]      Enable space level security group definitions
  --isolation-segment=                      Isolation segment assigned to space
  --named-asg=                              Named asg(s) to assign to space, specify multiple times
  --named-quota=                            Named quota to assign to space

quota:
  --enable-space-quota=[true|false]         Enable the Space Quota in the config
  --memory-limit=                           An Space's memory limit in Megabytes
  --instance-memory-limit=                  Space Application instance memory limit in Megabytes
  --total-routes=                           Total Routes capacity for an Space
  --total-services=                         Total Services capacity for an Space
  --paid-service-plans-allowed=[true|false] Allow paid services to appear in an Space
  --total-reserved-route-ports=             Total Reserved Route Ports capacity for an Space
  --total-service-keys=                     Total Service Keys capacity for an Space
  --app-instance-limit=                     App Instance Limit for a space
  --app-task-limit=                         App Task Limit for a space
  --log-rate-limit-bytes-per-second=        Log Rate limit per app for a space

developer:
  --developer-ldap-user=                    Ldap User to add, specify multiple times
  --developer-user=                         User to add, specify multiple times
  --developer-saml-user=                    SAML user to add, specify multiple times
  --developer-ldap-group=                   Group to add, specify multiple times

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
```

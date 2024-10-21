&larr; [back to Commands](../README.md)

# `cf-mgmt-config global`

`global` allows for updating cf-mgmt.yml global configuration values.

## Command Usage

```
Usage:
  cf-mgmt-config [OPTIONS] global [global-OPTIONS]

Help Options:
  -h, --help                                        Show this help message

[global command options]
    --config-dir=                                   Name of the config directory (default: config) [$CONFIG_DIR]
    --enable-delete-isolation-segments=[true|false] Enable removing isolation segments
    --enable-delete-shared-domains=[true|false]     Enable removing shared domains
    --enable-service-access=[true|false]            Enable managing service access
    --enable-unassign-security-groups=[true|false]  Enable unassigning security groups
    --skip-unassign-security-group-regex=           Skip unassigning security groups for names matching regex
    --metadata-prefix=                              Prefix for org/space metadata
    --enable-metadata-prefix=[true|false]           Enable useing metadata prefixes
    --staging-security-group=                       Staging Security Group to add
    --remove-staging-security-group=                Staging Security Group to remove
    --running-security-group=                       Running Security Group to add
    --remove-running-security-group=                Running Security Group to remove
    --shared-domain=                                Shared Domain to add
    --router-group-shared-domain=                   Router Group Shared Domain to add
    --router-group-shared-domain-group=             Router Group Shared Domain group
    --internal-shared-domain=                       Internal Shared Domain to add
    --remove-shared-domain=                         Shared Domain to remove

    service-access:
      --broker=                                     Name of Broker
      --service=                                    Name of Service
      --all-access-plan=                            Plan to give access to all orgs
      --limited-access-plan=                        Plan to give limited access to, must also provide org list
      --org=                                        Orgs to add to limited plan
      --remove-org=                                 Orgs to remove from limited plan
      --no-access-plan=                             Plan to ensure no access for any org
```

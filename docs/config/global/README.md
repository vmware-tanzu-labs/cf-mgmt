&larr; [back to Commands](../README.md)

# `cf-mgmt-config global`

`global` allows for updating cf-mgmt.yml global configuration values.

## Command Usage

```
Usage:
  cf-mgmt-config [OPTIONS] global [global-OPTIONS]

Help Options:
  -h, --help                                              Show this help message

[global command options]
  --config-dir=                                   Name of the config directory (default: config) [$CONFIG_DIR]
  --enable-delete-isolation-segments=[true|false] Enable removing isolation segments
  --enable-delete-shared-domains=[true|false]     Enable removing shared domains
  --enable-service-access=[true|false]            Enable managing service access
  --enable-unassign-security-groups=[true|false]  Enable unassigning security groups
  --metadata-prefix=                              Prefix for org/space metadata
  --staging-security-group=                       Staging Security Group to add
  --remove-staging-security-group=                Staging Security Group to remove
  --running-security-group=                       Running Security Group to add
  --remove-running-security-group=                Running Security Group to remove
  --shared-domain=                                Shared Domain to add
  --router-group-shared-domain=                   Router Group Shared Domain to add
  --router-group-shared-domain-group=             Router Group Shared Domain group
  --internal-shared-domain=                       Internal Shared Domain to add
  --remove-shared-domain=                         Shared Domain to remove
```

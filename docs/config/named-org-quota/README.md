&larr; [back to Commands](../README.md)

# `cf-mgmt-config named-org-quota`

`named-org-quota` command allows creating/updating named org quota

## Command Usage

```
Usage:
  cf-mgmt-config [OPTIONS] named-org-quota [named-org-quota-OPTIONS]

Help Options:
  -h, --help                                        Show this help message

[named-org-quota command options]
      --config-dir=                             Name of the config directory (default: config) [$CONFIG_DIR]
      --name=                                   Name of quota

quota:
      --memory-limit=                           An Org's memory limit in Megabytes
      --instance-memory-limit=                  Global Org Application instance memory limit in Megabytes
      --total-routes=                           Total Routes capacity for an Org
      --total-services=                         Total Services capacity for an Org
      --paid-service-plans-allowed=[true|false] Allow paid services to appear in an org
      --total-private-domains=                  Total Private Domain capacity for an Org
      --total-reserved-route-ports=             Total Reserved Route Ports capacity for an Org
      --total-service-keys=                     Total Service Keys capacity for an Org
      --app-instance-limit=                     App Instance Limit an Org
      --app-task-limit=                         App Task Limit an Org
      --log-rate-limit-bytes-per-second=        Log Rate limit per app for an org

```

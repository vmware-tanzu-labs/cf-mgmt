&larr; [back to Commands](../README.md)

# `cf-mgmt-config named-space-quota`

`named-space-quota` command allows creating/updating named space quotas for a given org

## Command Usage

```
Usage:
  cf-mgmt-config [OPTIONS] named-space-quota [named-space-quota-OPTIONS]

Help Options:
  -h, --help                                        Show this help message

[named-space-quota command options]
      --config-dir=                             Name of the config directory (default: config) [$CONFIG_DIR]
      --name=                                   Name of quota
      --org=                                    Name of org

quota:
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
```

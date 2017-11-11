&larr; [back to Commands](../README.md)

# `cf-mgmt init-config`

`init-config` will initialize a folder structure to add a ldap.yml and orgs.yml file.  This should be where you start to leverage cf-mgmt.  If your foundation is ldap enabled you can specify the ldap configuration info in ldap.yml otherwise you can disable this feature by setting the flag to false.

## Command Usage

```
Usage:
  main [OPTIONS] init-config [init-config-OPTIONS]

Help Options:
  -h, --help            Show this help message

[init-config command options]
  --config-dir= Name of the config directory (default: config) [$CONFIG_DIR]
```

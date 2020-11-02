&larr; [back to Commands](../README.md)

# `cf-mgmt-config update-orgs`

`update-orgs` command allows updating orgs.yml

_**Note**: If you intend to enable the deletion of orgs, please double check your `protected-orgs` list in your `orgs.yml`. Consider adding any orgs that might be created as part of a Tile installation._

## Command Usage
```
Usage:
  cf-mgmt-config [OPTIONS] update-orgs [update-orgs-OPTIONS]

Help Options:
  -h, --help                                Show this help message

[update-orgs command options]
  --config-dir=                     Name of the config directory (default: config) [$CONFIG_DIR]
  --enable-delete-orgs=[true|false] Enable delete orgs option
  --protected-org=                  Add org(s) to protected org list, specify multiple times. Uses re2 syntax:
                                    https://github.com/google/re2/wiki/Syntax
  --protected-org-to-remove=        Remove org(s) from protected org list, specify multiple times
```

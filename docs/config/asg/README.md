&larr; [back to Commands](../README.md)

# `cf-mgmt-config asg`

`asg` will add/update a <asg-name>.json file in the asgs folder of the configuration

## Command Usage

```
Usage:
  cf-mgmt-config [OPTIONS] asg [add-asg-OPTIONS]

Help Options:
  -h, --help            Show this help message

[asg command options]
  --config-dir= Name of the config directory (default: config) [$CONFIG_DIR]
  --asg=        ASG name
  --path=       path to asg definition file
  --override    override current definition
  --type=[space|default] Space asg or default asg (default: space)

```

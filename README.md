# cloud foundry mgmt
Go automation for managing orgs, spaces that can be driven from concourse pipeline and GIT managed metadata

## Install
- Either download a compiled release for your platform (make sure on linux/mac you chmod +x the binary)

or

- go get github.com/pivotalservices/cf-mgmt

### The following operation are enabled with cf-mgmt for helping to manage your configuration

- init-config

```
USAGE:
   cf-mgmt init-config [command options] [arguments...]

DESCRIPTION:
   initializes folder structure for configuration

OPTIONS:
   --config-dir value  config dir.  Default is config [$CONFIG_DIR]
```

- add-org-to-config

```
USAGE:
   cf-mgmt add-org-to-config [command options] [arguments...]

DESCRIPTION:
   adds specified org to configuration

OPTIONS:
   --org value         org name to add [$ORG]
   --config-dir value  config dir.  Default is config [$CONFIG_DIR]
```

- add-space-to-config            

```
USAGE:
   cf-mgmt add-space-to-config [command options] [arguments...]

DESCRIPTION:
   adds specified space to configuration for org

OPTIONS:
   --config-dir value  config dir.  Default is config [$CONFIG_DIR]
   --org value         org name of space [$ORG]
   --space value       space name to add [$SPACE]

```

- generate-concourse-pipeline

```
USAGE:
   cf-mgmt generate-concourse-pipeline [arguments...]

DESCRIPTION:
   generate-concourse-pipeline
```   

### The following operation are enabled with cf-mgmt that will leverage configuration to modify your Cloud Foundry installation

- create-orgs
`                   `
- update-org-quotas             
- update-org-users              
- create-spaces                 
- update-spaces                 
- update-space-quotas           
- update-space-users            
- update-space-security-groups  

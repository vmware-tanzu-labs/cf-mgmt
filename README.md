# cloud foundry mgmt
Go automation for managing orgs, spaces that can be driven from concourse pipeline and GIT managed metadata

## Install
Either download a compiled release for your platform (make sure on linux/mac you chmod +x the binary)

**or**

```
go get github.com/pivotalservices/cf-mgmt
```

## Contributing
PRs are always welcome or open issues if you are experiencing an issue and will do my best to address issues in timely fashion.

### The following operation are enabled with cf-mgmt for helping to manage your configuration

#### init-config

This command will initialize a folder structure to add a ldap.yml and orgs.yml file.  This should be where you start to leverage cf-mgmt.  If your foundation is ldap enabled you can specify the ldap configuration info in ldap.yml otherwise you can disable this feature by setting the flag to false.

```
USAGE:
   cf-mgmt init-config [command options] [arguments...]

DESCRIPTION:
   initializes folder structure for configuration

OPTIONS:
   --config-dir value  config dir.  Default is config [$CONFIG_DIR]
```

#### add-org-to-config

This will add the specified org to orgs.yml and create folder based on the org name you specified.  Within this folder will be an orgConfig.yml and spaces.yml which will be empty but will eventually contain a list of spaces.  Any org listed in orgs.yml will be created when the create-orgs operation is ran.

orgConfig.yml allows specifying for the following:
- what groups to map to org roles (OrgManager, OrgBillingManager, OrgAuditor)
- setting up quotas for the org

```
USAGE:
   cf-mgmt add-org-to-config [command options] [arguments...]

DESCRIPTION:
   adds specified org to configuration

OPTIONS:
   --org value         org name to add [$ORG]
   --config-dir value  config dir.  Default is config [$CONFIG_DIR]
```

#### add-space-to-config

This command allows for adding a space to a previously defined org.  This will generate a folder for each space inside the orgs folder.  In the spaces folder will contain a spaceConfig.yml and a security-group.json file.  Any space listed in spaces.yml will be created when the create-spaces operation is ran.  The spaceConfig.yml allows for specifying the following:   

- allow ssh at space level
- map ldap group names to SpaceDeveloper, SpaceManager, SpaceAuditor role
- setup quotas at a space level (if enabled)
- apply application security group config at space level (if enabled)        

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

#### generate-concourse-pipeline

This will generate a pipeline.yml, vars.yml and necessary task yml files for running all the tasks listed below.  Just need to update your vars.yml check in all your code to GIT and execute the fly command to register your pipeline.  

```
USAGE:
   cf-mgmt generate-concourse-pipeline [arguments...]

DESCRIPTION:
   generate-concourse-pipeline
```   

### The following operation are enabled with cf-mgmt that will leverage configuration to modify your Cloud Foundry installation

To execute any of the following you will need to provide:
- **user id** that has priviledges to create orgs/spaces
- **password** for the above user account
- **uaac client secret** for the account that can add users (assumes the same user account for cf commands is used)
- **system domain** name of your foundation

#### create-orgs
- creates orgs specified in orgs.yml

#### update-org-quotas
- updates org quotas specified in orgConfig.yml

#### update-org-users              
- syncs users from ldap groups configured in orgConfig.yml assuming that ldap.yml is configured

#### create-spaces                 
- creates spaces for all spaces listed in each spaces.yml

#### update-spaces                 
- updates allow ssh into space property

#### update-space-quotas           
- creates/updates quota for a given space

#### update-space-users
- syncs users from ldap groups configured in spaceConfig.yml assuming that ldap.yml is configured

#### update-space-security-groups  
- creates/updates application security groups for a given space

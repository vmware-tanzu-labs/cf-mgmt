# cloud foundry mgmt
Go automation for managing orgs, spaces that can be driven from concourse pipeline and GIT managed metadata


## Install
Either download a compiled release for your platform (make sure on linux/mac you chmod +x the binary)

**or**

```
go get github.com/pivotalservices/cf-mgmt
```

## Testing

Get ready for tests:
```
go get github.com/onsi/ginkgo
go get github.com/onsi/gomega
go get github.com/golang/mock/gomock
go get github.com/golang/mock/mockgen
go get github.com/golang/protobuf/proto
./update-mocks.sh
```
Run tests:
```
docker pull cwashburn/ldap
docker run -d -p 389:389 --name ldap -t cwashburn/ldap
go test $(glide nv) -v
```

## Wercker cli tests
```
./testrunner
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

### Configuration
After running the above commands there will be a config directory in the working directory.  This will have a folder per org and within each org there will be a foler for each space.

```
├── ldap.yml
├── orgs.yml
├── test
│   ├── orgConfig.yml
│   ├── space1
│   │   ├── security-group.json
│   │   └── spaceConfig.yml
│   └── spaces.yml
└── test2
    ├── orgConfig.yml
    ├── space1a
    │   ├── security-group.json
    │   └── spaceConfig.yml
    └── spaces.yml
```

#### Org Configuration
There is a orgs.yml that contains list of orgs that will be created.  This should have a cooresponding folder with name of the orgs cf-mgmt is managing.  This will contain a orgConfig.yml and folder for each space.  Each orgConfig.yml consists of the following.

```
# org name
org: test

org-billingmanager:
  # list of ldap users that will be created in cf and given billing manager role
  ldap_users:
    - cwashburn1
    - cwashburn2

  # list of users that would be given billing manager role (must already be a user created via cf create-user)
  users:
    - cwashburn@testdomain.com
    - cwashburn2@testdomain.com


  # ldap group that contains users that will be added to cf and given billing manager role
  ldap_group: test_billing_managers

org-manager:
  # list of ldap users that will be created in cf and given org manager role
  ldap_users:
    - cwashburn1
    - cwashburn2

  # list of users that would be given org manager role (must already be a user created via cf create-user)
  users:
    - cwashburn@testdomain.com
    - cwashburn2@testdomain.com

  # ldap group that contains users that will be added to cf and given org manager role
  ldap_group: test_org_managers

org-auditor:
  # list of ldap users that will be created in cf and given org manager role
  ldap_users:
    - cwashburn1
    - cwashburn2

  # list of users that would be given org auditor role (must already be a user created via cf create-user)
  users:
    - cwashburn@testdomain.com
    - cwashburn2@testdomain.com

  # ldap group that contains users that will be added to cf and given org auditor role
  ldap_group: test_org_auditors

# if you wish to enable custom org quotas
enable-org-quota: true
# 10 GB limit
memory-limit: 10240
# unlimited
instance-memory-limit: -1
total-routes: 10
# unlimited
total-services: -1
paid-service-plans-allowed: true
```

#### Space Configuration
There will be a spaces.yml that will list all the spaces for each org.  There will also be a folder for each space with the same name.  Each folder will contain a spaceConfig.yml and security-group.json file with an empty json file.  Each spaceConfig.yml will have the following configuration options.  

```
---
# org that is space belongs to
org: test

# space name
space: space1

# if cf ssh is allowed for space
allow-ssh: yes

space-manager:
  # list of ldap users that will be created in cf and given space manager role
  ldap_users:
    - cwashburn1
    - cwashburn2

  # list of users that would be given space manager role (must already be a user created via cf create-user)
  users:
    - cwashburn@testdomain.com
    - cwashburn2@testdomain.com

  # ldap group that contains users that will be added to cf and given space manager role
  ldap_group: test_space1_managers

space-auditor:
  # list of ldap users that will be created in cf and given space auditor role
  ldap_users:
    - cwashburn1
    - cwashburn2

  # list of users that would be given space auditor role (must already be a user created via cf create-user)
  users:
    - cwashburn@testdomain.com
    - cwashburn2@testdomain.com

  # ldap group that contains users that will be added to cf and given space auditor role
  ldap_group: test_space1_auditors

space-developer:
  # list of ldap users that will be created in cf and given space developer role
  ldap_users:
    - cwashburn1
    - cwashburn2

  # list of users that would be given space developer role (must already be a user created via cf create-user)
  users:
    - cwashburn@testdomain.com
    - cwashburn2@testdomain.com

  # ldap group that contains users that will be added to cf and given space developer role
  ldap_group: test_space1_developers

# to enable custom quota at space level  
enable-space-quota: true
# 10 GB limit
memory-limit: 10240
# unlimited
instance-memory-limit: -1
total-routes: 10
# unlimited
total-services: -1
paid-service-plans-allowed: true

# to enable custom asg for the space.  If true will deploy asg defined in security-group.json within space folder
enable-security-group: false
```


### Recommended workflow

Operations team can setup a a git repo seeded with cf-mgmt configuration.  This will be linked to a concourse pipeline (example pipeline generated below) that will create orgs, spaces, map users, create quotas, deploy ASGs based on changes to git repo.  Consumers of this can submit a pull request via GIT to the ops team with comments like any other commit.  This will create a complete audit log of who requested this and who approved within GIT history.  Once PR accepted then concourse will provision the new items.

#### generate-concourse-pipeline

This will generate a pipeline.yml, vars.yml and necessary task yml files for running all the tasks listed below.  Just need to update your vars.yml check in all your code to GIT and execute the fly command to register your pipeline.  

```
USAGE:
   cf-mgmt generate-concourse-pipeline [arguments...]

DESCRIPTION:
   generate-concourse-pipeline
```   

### Known Issues
Currently does not remove anything that is not in configuration.  All functions are additive.  So removing users, orgs, spaces is not currently a function if they are not configured in cf-mgmt but future plans are to have a flag to opt-in for this feature.  This will likely start with removing users that are not configured in the orgs/spaces managed by cf-mgmt.

### The following operation are enabled with cf-mgmt that will leverage configuration to modify your Cloud Foundry installation

To execute any of the following you will need to provide:
- **user id** that has privileges to create orgs/spaces
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

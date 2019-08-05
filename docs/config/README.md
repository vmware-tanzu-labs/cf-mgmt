# DEPRECATED Commands
* [add-org](add-org/README.md), use [org](org/README.md)
* [add-space](add-space/README.md), use [space](space/README.md)
* [add-asg](add-asg/README.md), use [asg](asg/README.md)
* [update-org](update-org/README.md), use [org](org/README.md)
* [update-space](update-space/README.md), use [space](space/README.md)

# Configuration Commands
* [init](init/README.md)
* [global](global/README.md)
* [asg](asg/README.md)
* [org](org/README.md)
* [space](space/README.md)
* [delete-org](delete-org/README.md)
* [delete-space](delete-space/README.md)
* [generate-concourse-pipeline](generate-concourse-pipeline/README.md)
* [update-orgs](update-orgs/README.md)
* [rename-org](rename-org/README.md)
* [rename-space](rename-space/README.md)
* [named-org-quota](named-org-quota/README.md)
* [named-space-quota](named-space-quota/README.md)
* [version ](version/README.md)

### Global Config
There is global configuration that is managed in `cf-mgmt.yml`.  The following options exist in that configuration.

```
enable-delete-isolation-segments: false #true/false
enable-unassign-security-groups: false #true/false
running-security-groups: # array of security groups to apply to running
- all_access
- public_networks
- dns
- load_balancer
staging-security-groups: # array of security groups to apply to staging
- all_access
- public_networks
- dns
shared-domains: # map of shared domains and their configuration 1.0.12+
  dev.cfdev.sh: #shared domain name
    internal: false
  dev.cfdev.sh.tcp: #shared domain name
    internal: false
    router-group: default-tcp #router group to associate with domain
enable-remove-shared-domains: true #true/false

enable-service-access: true #true/false

# added in v1.0.31
service-access:
- broker: dedicated-mysql-broker
  services:
  - service: p.mysql
    all_access_plans: # db-small plan is available to all orgs
    - db-small  
    limited_access_plans:# db-medium is only available to cfdev-org as well as any org that is in `protected` orgs list
    - plan: db-medium
      orgs:
      - cfdev-org
    no_access_plans: # disables db-large plan for all orgs
    - db-large
- broker: rabbitmq-odb
  services:
  - service: p.rabbitmq
    all_access_plans:
    - single-node-3.7
- broker: p-rabbitmq
  services:
  - service: p-rabbitmq
    all_access_plans:
    - standard
```

#### Org Configuration
There is a orgs.yml that contains list of orgs that will be created.  This should have a corresponding folder with name of the orgs cf-mgmt is managing. orgs.yml also can be configured with a list of protected orgs which would never be deleted when using the the `delete-orgs` command. An example of how orgs.yml could be configured is seen below.

```
orgs:
- foo-org
- bar-org
# added in 0.0.63+ which will remove orgs not configured in cf-mgmt
enable-delete-orgs: true
# added in 0.0.63+ which allows configuration of orgs to 'ignore'
protected_orgs:
- system
- p-spring-cloud-services

```

This will contain a orgConfig.yml and folder for each space.  Each orgConfig.yml consists of the following.
```
# org name
org: test

# added in 1.0.9+ to allow renaming orgs
original-org: foo

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

  # added in 0.0.62+ which will allow configuration of a list of groups works with ldap_group
  ldap_groups:
    - test_billing_managers_2

  # added in 0.0.66+ which will allow configuration of a list of saml user email addresses
  saml_users:
    - cwashburn@testdomain.com
    - cwashburn2@testdomain.com
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

  # added in 0.0.62+ which will allow configuration of a list of groups works with ldap_group
  ldap_groups:
    - test_org_managers_2

  # added in 0.0.66+ which will allow configuration of a list of saml user email addresses
  saml_users:
    - cwashburn@testdomain.com
    - cwashburn2@testdomain.com
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

  # added in 0.0.62+ which will allow configuration of a list of groups works with ldap_group
  ldap_groups:
    - test_org_auditors_2

  # added in 0.0.66+ which will allow configuration of a list of saml user email addresses
  saml_users:
    - cwashburn@testdomain.com
    - cwashburn2@testdomain.com
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

# added in 0.0.48+ which will remove users from roles if not configured in cf-mgmt
enable-remove-users: true/false

# added in 0.0.64+ which will remove users from roles if not configured in cf-mgmt
private-domains: ["test.com", "test2.com"]
enable-remove-private-domains: true/false

# added in 1.0.12+ allows specifying a named quota, cannot be used with enable-org-quota
named-quota:

# added in 1.0.26+ allows adding metadata to orgs and spaces (requires cf v3 3.66.0 or greater)
metadata:
  labels:
    foo: bar
  annotations:
    hello: world
```

#### Space Configuration
There will be a spaces.yml that will list all the spaces for each org.  There will also be a folder for each space with the same name.  Each folder will contain a spaceConfig.yml and security-group.json file with an empty json file.  

Each spaceConfig.yml will have the following configuration options:
- allow ssh at space level
- map ldap group names to SpaceDeveloper, SpaceManager, SpaceAuditor role
- setup quotas at a space level (if enabled)
- apply application security group config at space level (if enabled)    

```
# org that is space belongs to
org: test

# space name
space: space1

# added in 1.0.9+ to allow renaming spaces
original-space: old-space1

# if cf ssh is allowed for space
allow-ssh: yes

# to temporarily grant ssh access added in 1.0.13+, use cf-mgmt-config to specify time as field needs to be in RFC3339 format
allow-ssh-until: "2019-01-13T18:09:16-07:00"

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

  # added in 0.0.62+ which will allow configuration of a list of groups works with ldap_group
  ldap_groups:
    - test_space1_managers_2

  # added in 0.0.66+ which will allow configuration of a list of saml user email addresses
  saml_users:
    - cwashburn@testdomain.com
    - cwashburn2@testdomain.com
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

  # added in 0.0.62+ which will allow configuration of a list of groups works with ldap_group
  ldap_groups:
    - test_space1_auditors_2

  # added in 0.0.66+ which will allow configuration of a list of saml user email addresses
  saml_users:
    - cwashburn@testdomain.com
    - cwashburn2@testdomain.com

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

  # added in 0.0.62+ which will allow configuration of a list of groups works with ldap_group
  ldap_groups:
    - test_space1_developers_2

  # added in 0.0.66+ which will allow configuration of a list of saml user email addresses
  saml_users:
    - cwashburn@testdomain.com
    - cwashburn2@testdomain.com
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

# added in 0.0.48+ which will remove users from roles if not configured in cf-mgmt
enable-remove-users: true/false

# allowing security groups to be applied that are defined globally
named-security-groups: []

# added in 1.0.12+ allows specifying a named quota, cannot be used with enable-space-quota
named-quota:

# added in 1.0.26+ allows unassigning named security groups that are not in configuration
enable-unassign-security-group: true/false

# added in 1.0.26+ allows adding metadata to orgs and spaces (requires cf v3 3.66.0 or greater - PCF 2.5+)
metadata:
  labels:
    foo: bar
  annotations:
    hello: world
```

#### Space Default Configuration

The file spaceDefaults.yml can be used to specify a default set of roles for user and groups to be applied to all spaces.  
This will be merged with the space-specific roles.  
Note that this is actually processed at runtime, not when spaces are added to the config.  

### LDAP Configuration
LDAP configuration file ```ldap.yml``` is located under the ```config``` folder. By default, LDAP is disabled and you can enable it by setting ```enabled: true```. Once this is enabled, all other LDAP configuration properties are required.

```
enabled: true
ldapHost: 127.0.0.1
ldapPort: 10389
#true/false (default false)
use_tls: true
bindDN: uid=admin,ou=system
userSearchBase: ou=users,dc=example,dc=com
userNameAttribute: uid
# optional added in v1.0.20+
userObjectClass: <object class that matches your ldap/active directory configuration for users (inetOrgPerson, organizationalPerson)>
userMailAttribute: mail
groupSearchBase: ou=groups,dc=example,dc=com
groupAttribute: member
# optional added in v1.0.20+
groupObjectClass: <object class that matches your ldap/active directory configuration for groups (group, groupOfNames)>
origin: ldap

# optional added in 1.0.11+ - true/false
insecure_skip_verify: false
# optional added in 1.0.11+ if ldap server is signed by non-public CA provide ca pem here
ca_cert: |
```

### SAML Configuration with ldap group lookups
LDAP configuration file ```ldap.yml``` is located under the ```config``` folder. To have cf-mgmt create SAML users in UAA need to enable ldap to lookup the user information from an LDAP source to properly create the SAML users.  In orgConfig.yml and spaceConfig.yml leverage either/or `ldap_users` or `ldap_group(s)`  

```
enabled: true
ldapHost: 127.0.0.1
ldapPort: 10389
#true/false (default false)
use_tls: true
bindDN: uid=admin,ou=system
userSearchBase: ou=users,dc=example,dc=com
userNameAttribute: uid
# optional added in v1.0.20+
userObjectClass: <object class that matches your ldap/active directory configuration for users (inetOrgPerson, organizationalPerson)>
userMailAttribute: mail
groupSearchBase: ou=groups,dc=example,dc=com
groupAttribute: member
# optional added in v1.0.20+
groupObjectClass: <object class that matches your ldap/active directory configuration for groups (group, groupOfNames)>
origin: <needs to match origin configured for elastic runtime>

# optional added in 1.0.11+ - true/false
insecure_skip_verify: false
# optional added in 1.0.11+ if ldap server is signed by non-public CA provide ca pem here
ca_cert: |
```

### SAML Configuration
LDAP configuration file ```ldap.yml``` is located under the ```config``` folder. To have cf-mgmt create SAML users you can disable ldap integration for looking up users in ldap groups with v0.0.66+ as orgConfig.yml and spaceConfig.yml now includes a saml_users array attribute which can contain a list of email addresses.

```
enabled: false
origin: <needs to match origin configured for elastic runtime>
ldapHost:
ldapPort: 389
bindDN:
userSearchBase:
userNameAttribute:
userMailAttribute:
groupSearchBase:
groupAttribute:
```

### Enable Temporary Application SSH Access
With 1.0.13+ there is ability to grant applicaiton ssh access for a specific duration.  Durations supported are in number of Days (D), Hours (H) or Minutes (M).  Use the cf-mgmt-config cli to update a given space with one of these metrics.  This will generate the timestamp in the correct format for you.  You must also use the latest generated concourse pipeline as this places update-space command on a timer to run every 15m (by default) to check to see if time has elapsed to re-disable application ssh access

The following will enable for 2 days:
```
cf-mgmt-config update-space --config-dir <your directory> --org <org> --space <space> --allow-ssh false --allow-ssh-until 2D  
```

The following will enable for 5 hours:
```
cf-mgmt-config update-space --config-dir <your directory> --org <org> --space <space> --allow-ssh false --allow-ssh-until 5H
```

The following will enable for 95 minutes:
```
cf-mgmt-config update-space --config-dir <your directory> --org <org> --space <space> --allow-ssh false --allow-ssh-until 95M
```

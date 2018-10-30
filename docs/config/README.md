# Configuration Commands
* [init](init/README.md)
* [add-org](add-org/README.md)
* [add-space](add-space/README.md)
* [add-asg](add-asg/README.md)
* [delete-org](delete-org/README.md)
* [delete-space](delete-space/README.md)
* [generate-concourse-pipeline](generate-concourse-pipeline/README.md)
* [update-org](update-org/README.md)
* [update-orgs](update-orgs/README.md)
* [update-space](update-space/README.md)
* [rename-org](rename-org/README.md)
* [rename-space](rename-space/README.md)
* [version ](version/README.md)

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

# added in 1.0.9+ which allows configuration of service access
service-access:
  p-mysql: ["small","large"]
  p-rabbit: ["*"]
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
```

#### Space Default Configuration

The file spaceDefaults.yml can be used to specify a default set of roles for user and groups to be applied to all spaces.  
This will be merged with the space-specific roles.  
Note that this is actually processed at runtime, not when spaces are added to the config.  

### LDAP Configuration
LDAP configuration file ```ldap.yml``` is located under the ```config``` folder. By default, LDAP is disabled and you can enable it by setting ```enabled: true```. Once this is enabled, all other LDAP configuration properties are required.

### SAML Configuration with ldap group lookups
LDAP configuration file ```ldap.yml``` is located under the ```config``` folder. To have cf-mgmt create SAML users in UAA need to enable ldap to lookup the user information from an LDAP source to properly create the SAML users.  In orgConfig.yml and spaceConfig.yml leverage either/or `ldap_users` or `ldap_group(s)`  

```
enabled: true
ldapHost: 127.0.0.1
ldapPort: 10389
bindDN: uid=admin,ou=system
userSearchBase: ou=users,dc=example,dc=com
userNameAttribute: uid
userMailAttribute: mail
groupSearchBase: ou=groups,dc=example,dc=com
groupAttribute: member
origin: <needs to match origin configured for elastic runtime>
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

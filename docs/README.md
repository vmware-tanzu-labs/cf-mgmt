# Configuration Commands
**(DEPRECATED use cf-mgmt-config instead)**
* [add-org-to-config](add-org-to-config/README.md)
* [add-space-to-config](add-space-to-config/README.md)
* [generate-concourse-pipeline](generate-concourse-pipeline/README.md)
* [init-config](init-config/README.md)


# Commands

The following commands are enabled with cf-mgmt that will leverage configuration interact with your Cloud Foundry installation and need to authenticate

## Authentication requirements

Introduced in cf-mgmt 0.0.66+ is the ability to definee a non admin uaa client.  With this release the password field has been deprecated. To create a non-admin client execute the following command with [Cloud Foundry UAA Client](https://docs.pivotal.io/pivotalcf/1-11/adminguide/uaa-user-management.html).  

```
$ uaac target uaa.<system-domain>
$ uaac token client get <adminuserid> -s <admin-client-secret>

$ uaac client add cf-mgmt \
  --name cf-mgmt \
  --secret <cf-mgmt-secret> \
  --authorized_grant_types client_credentials,refresh_token \
  --authorities cloud_controller.admin,scim.read,scim.write
```

As you can see, `cloud_controller.admin,scim.read,scim.write` gives this user just enough rights to add/update/delete users, orgs and space and still being a non-admin user. Learn more about the scopes authorized by UAA at [UAA Scopes](https://github.com/cloudfoundry/uaa/blob/master/docs/UAA-APIs.rst#scopes-authorized-by-the-uaa)


To execute any of the following you will need to provide:
- **user-id** that has privileges to create/update/delete users, orgs and spaces. This user doesn't have to be an admin user. Assuming you have [Cloud Foundry UAA
- **client-secret** for the above user (assumes the same user account for cf commands is used)
- **system-domain** name of your foundation

Prior to v0.0.66 a **password** was also needed as you had to provide both a uaa user and uaa client.  This field has been deprecated and will be removed in a future release as going forward cf-mgmt will require a uaa client per the authentication directions.

* [create-org-private-domains](create-org-private-domains/README.md)
* [share-org-private-domains](share-org-private-domains/README.md)
* [create-orgs](create-orgs/README.md)
* [create-security-groups](create-security-groups/README.md)
* [assign-default-security-groups](assign-default-security-groups/README.md)
* [create-spaces](create-spaces/README.md)
* [delete-orgs](delete-orgs/README.md)
* [delete-spaces](delete-spaces/README.md)
* [export-config](export-config/README.md)
* [isolation-segments](isolation-segments/README.md)
* [update-org-quotas](update-org-quotas/README.md)
* [update-org-users](update-org-users/README.md)
* [update-space-quotas](update-space-quotas/README.md)
* [update-space-security-groups](update-space-security-groups/README.md)
* [update-space-users](update-space-users/README.md)
* [update-spaces](update-spaces/README.md)
* [version](version/README.md)

# Features
- Removing users from cf that are not in cf-mgmt metadata was added in 0.48+ release.  This is an opt-in feature for existing cf-mgmt users at an org and space config level.  For any new orgs/config created with cf-mgmt cli 0.48+ it will default this parameter to true.  To opt-in ensure you are using latest cf-mgmt version when running pipeline and add `enable-remove-users: true` to your configuration.

- Removing orgs and spaces from cf that are not in cf-mgmt metadata was added in 0.0.63+ release.  This is an opt-in feature for existing cf-mgmt users at an org and space config level.  For any new orgs/config created with cf-mgmt cli 0.0.63+ it will default this parameter to true.  To opt-in ensure you are using latest cf-mgmt version when running pipeline and add `enable-delete-orgs: true` or `enable-delete-spaces: true` to your configuration.

- Managing private domains at org level was added with 0.0.64+.  This requires you to update concourse pipeline to to invoke `create-org-private-domains` command.  By default `enable-remove-private-domains: true` is set for any new orgs created with 0.0.64+ cli.  This will remove any private domains for that org that are not in array of private domain names.

# Recommended workflow

Operations team can setup a a git repo seeded with cf-mgmt configuration.  This will be linked to a concourse pipeline (example pipeline generated below) that will create orgs, spaces, map users, create quotas, deploy ASGs based on changes to git repo.  Consumers of this can submit a pull request via GIT to the ops team with comments like any other commit.  This will create a complete audit log of who requested this and who approved within GIT history.  Once PR accepted then concourse will provision the new items.

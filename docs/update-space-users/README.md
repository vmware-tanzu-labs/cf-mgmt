&larr; [back to Commands](../README.md)

# `cf-mgmt update-space-quotas-users`

`update-space-users` command will:

- adds users from `ldap_groups` configured in spaceConfig.yml assuming that ldap.yml is configured
- adds ldap users in `ldap_users` configured in spaceConfig.yml assuming that ldap.yml is configured
- add internal `users` configured in spaceConfig.yml (internal users must exist in uaa first)
- add `saml_users` configured in spaceConfig.yml (internal users must exist in uaa first)
- will remove users from roles if `enable-remove-users` is set to `true` in spaceConfig.yml

## Command Usage

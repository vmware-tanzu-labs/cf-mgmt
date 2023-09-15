TAS version | Compatible?
--- | ---
4.0 | ![CI](https://runway-ci.eng.vmware.com/api/v1/teams/cryogenics/pipelines/cf-mgmt/jobs/test-against-tas-4-0/badge)
3.0 | ![CI](https://runway-ci.eng.vmware.com/api/v1/teams/cryogenics/pipelines/cf-mgmt/jobs/test-against-tas-3-0/badge)
2.13 | ![CI](https://runway-ci.eng.vmware.com/api/v1/teams/cryogenics/pipelines/cf-mgmt/jobs/test-against-tas-2-13/badge)
2.11 | ![CI](https://runway-ci.eng.vmware.com/api/v1/teams/cryogenics/pipelines/cf-mgmt/jobs/test-against-tas-2-11/badge)

# Cloud Foundry Management (cf-mgmt)

Go automation for managing orgs, spaces, users (from ldap groups or internal store) mapping to roles, quotas, application security groups and private-domains that can be driven from concourse pipeline and GIT managed metadata

## New Major Release Information

There has been major refactoring to internals of cf-mgmt to remove duplicate code that is not supported by go-cfclient library.  This release SHOULD be backward compatible but wanting to make community aware of a major change.  This will be released as the latest tag on dockerhub.  If you experience any problems you can revert your cf-mgmt to use the previously released version with tag `0.0.91`.

This can be done by modifying you cf-mgmt.yml concourse task with the following:

```yml
---
platform: linux

image_resource:
  type: docker-image
  source: {repository: pivotalservices/cf-mgmt, tag: "0.0.91"}

inputs:
  - name: config-repo

run:
  path: config-repo/ci/tasks/cf-mgmt.sh
```

## Overview

The __cf-mgmt__ tool is composed by 2 CLIs, `cf-mgmt` and `cf-mgmt-config`, providing the features to declarativelly manage orgs, spaces, users mapping to roles, quotas, application security groups and private-domains.

- __cf-mgmt-config__ CLI is resposible for creating the configuration files that represent the desired state of your foundation and provides the set of commands for you to manage this configuration.

- __cf-mgmt__ CLI is resposible to apply the configuration generated by `cf-mgmt-config` tool to the foundation. It provides commands to apply the config as a whole or only parts of it.

### Concourse

A common use of `cf-mgmt` is to generate a Concourse pipeline that apply the configuration generated by `cf-mgmt-config` to a foundation. This is done by using a git repository as a resouce for the configuration and having the pipeline to read from there and apply the changes. `cf-mgmt` provides a command to generate this pipeline. See more at the [Gettting Started](#getting-started) section below.

### Export Configuration

If there's already a configured foundation that you want to start using cf-mgmt on, `cf-mgmt export-config` command will export the current foundation configs and generate the files for cf-mgmt usage. See more in the [docs](docs/export-config/README.md).

## Getting Started

### Install

Compiled [releases](https://github.com/vmwarepivotallabs/cf-mgmt/releases) are available on Github.
Download the binary for your platform and place it somewhere on your path.
Don't forget to `chmod +x` the file on Linux and macOS.

### Create UAA Client

cf-mgmt needs a uaa client to be able to interact with cloud controller and uaa for create, updating, deleting, and listing entities.

To create a non-admin client execute the following command with [Cloud Foundry UAA Client](https://github.com/cloudfoundry/cf-uaac).  Recent addition of 2 authorities needed to setup shared domains with tcp routing `routing.router_groups.read`

```sh
uaac target uaa.<system-domain>
uaac token client get <adminuserid> -s <admin-client-secret>

uaac client add cf-mgmt \
  --name cf-mgmt \
  --secret <cf-mgmt-secret> \
  --authorized_grant_types client_credentials,refresh_token \
  --authorities cloud_controller.admin,scim.read,scim.write,routing.router_groups.read
```

Or with the [golang-based UAA CLI](https://github.com/cloudfoundry-incubator/uaa-cli):

`go install github.com/cloudfoundry-incubator/uaa-cli`

```sh
uaa-cli target https://uaa.<system-domain>

uaa-cli get-client-credentials-token <adminuserid> -s <admin-client-secret>

uaa-cli create-client cf-mgmt \
  --client_secret <cf-mgmt-secret> \
  --authorized_grant_types client_credentials,refresh_token \
  --authorities cloud_controller.admin,scim.read,scim.write,routing.router_groups.read
```

### Setup Configuration

Navigate into a directory in which will become your git repository for cf-mgmt configuration

1. Initialize git repository by either cloning a remote or using `git init`

    - This git repository is going to be used as a place to store the config files and will be consumed by a Concourse Pipeline. You should not push your `vars.yml` file or any other files with secrets to this repo.

2. Create initial configuration files. You can either setup your configuration by using `init` command or `export-config` if you want to start managing a foundation that already have some workspace setup:

   - [init](docs/config/init/README.md) command from `cf-mgmt-config` if you are wanting to start with a blank configuration and add the config using `cf-mgmt-config` operations
   - [export-config](docs/export-config/README.md) command from `cf-mgmt` if you have an existing foundation you can use this to reverse engineer your configuration.

> By default, `cf-mgmt` has the `enable-delete-orgs` option set to `false` to avoid unintentional deletions. If you'd like to have `cf-mgmt` handle deletion of orgs as well, please double check the list of `protected-orgs` and update the `enable-delete-orgs` flag to true in the `orgs.yml` of your config repository.
> You can also modify these settings using the [`cf-mgmt-config update-orgs` command](docs/config/update-orgs/README.md).

> Check the [config docs](docs/config/README.md#global-config) to understand the configuration files structure

3. *(optional)* Configure LDAP/SAML Options. If your foundation uses LDAP and/or SAML, you will need to configure ldap.yml with the correct values.

   - [LDAP only config](docs/config/README.md#ldap-configuration)
   - [SAML with LDAP groups](docs/config/README.md#saml-configuration-with-ldap-group-lookups)
   - [SAML only](docs/config/README.md#saml-configuration)

4. [Generate the concourse pipeline](docs/config/generate-concourse-pipeline/README.md) using `cf-mgmt-config`
    - ```cf-mgmt-config [OPTIONS] generate-concourse-pipeline [generate-concourse-pipeline-OPTIONS]```

5. Make sure you __.gitingore the vars.yml__ file that is generated: `echo vars.yml >> .gitignore`

6. Update your `vars.yml` file with your config git repo info, domains and the UAA client info you created in the [previous section](#create-uaa-client).
    - Use the UAA Client name as `user_id`

7. Commit and push your changes to your git repository.
    - After you `fly` the pipeline in the next step, the pipeline will observe this repository and start a new run everytime you push new configurations to the repository.

8. fly your pipeline after you have filled in vars.yml

```sh
fly -t <targetname> login -c <concourse_instance>
fly -t <targetname> set-pipeline -p <pipeline_name> \
   -c pipeline.yml \
   -l vars.yml \
   —-var "ldap_password=<ldap_password>" \
   --var "client_secret=<client_secret>" \
   —-var "password=<org/space_admin_password>"
```

>Check the [Concourse docs](https://concourse-ci.org/fly.html) if not familiar with the `fly` CLI

9. Look for your new Pipeline in your Concourse console. If everything was properly configured all tasks of your pipeline should execute successfully.

>You now have a Pipeline ready to apply configuration changes to your foundation. Explore the [docs](docs/config/README.md) to learn the available commands in `cf-mgmt-config`, try creating new workspace configs, then commit and push the files to you git repository. Your pipeline should kick in and apply the changes.

## Support

cf-mgmt is a community supported cloud foundry add-on.  Opening issues for questions, feature requests and/or bugs is the best path to getting "support".  We strive to be active in keeping this tool working and meeting your needs in a timely fashion.

### Install Binary

Compiled releases are available on Github.
Download the binary for your platform and place it somewhere on your path.
Don't forget to `chmod +x` the file on Linux and macOS.

Alternatively, you may wish to build from source.

### Debug Output

When opening an issue please provide debug level output (scrubbed for any customer info) by using latest generated pipeline and setting LOG_LEVEL: debug or modifying current pipeline if you are not using latest pipeline to add the following to specific job step params

```yml
params:
  LOG_LEVEL: debug
  ... existing params
```

## Development

### Build from the source

`cf-mgmt` is written in [Go](https://golang.org/).
To build the binary yourself, follow these steps:

- Install `Go`.
- Clone the repo
- Build:
  - `cd cf-mgmt`
  - `go build -o cf-mgmt cmd/cf-mgmt/main.go`
  - `go build -o cf-mgmt-config cmd/cf-mgmt-config/main.go`

To cross compile, set the `$GOOS` and `$GOARCH` environment variables.
For example: `GOOS=linux GOARCH=amd64 go build`.

### Testing

To run the unit tests, use `go test ./...`.

### SSH configuration
In order to use the key `allow-ssh-until` in your space config, you must set
your `allow-ssh` to false. `cf-mgmt` treats a null value differently than false.

#### Integration tests

There are integration tests that require some additional configuration.

The LDAP tests require an LDAP server, which can be started with Docker:

```docker
docker pull cwashburn/ldap
docker run -d -p 389:389 --name ldap -t cwashburn/ldap
RUN_LDAP_TESTS=true go test ./ldap_integration/...
```

The remaining integration tests require [PCF Dev](https://pivotal.io/pcf-dev)
to be running, the CF CLI, and the [UAA CLI](https://github.com/cloudfoundry-incubator/uaa-cli).

```sh
cf dev start
uaa-cli target https://uaa.dev.cfdev.sh -k

uaa-cli get-client-credentials-token admin -s admin-client-secret

uaa-cli create-client cf-mgmt \
  --client_secret cf-mgmt-secret \
  --authorized_grant_types client_credentials,refresh_token \
  --authorities cloud_controller.admin,scim.read,scim.write,routing.router_groups.read
RUN_INTEGRATION_TESTS=true go test ./integration/...
```

### Code Generation

Some portions of this code are autogenerated.
To regenerate run `go generate ./...` from the project directory, or `go generate .` from a specific directory.

## Contributing

PRs are always welcome or open issues if you are experiencing an issue and will do my best to address issues in timely fashion.

## Documentation

- See [here](docs/README.md) for documentation on all the available commands for running cf-mgmt
- See [here](docs/config/README.md) for documentation on all the configuration documentation and commands

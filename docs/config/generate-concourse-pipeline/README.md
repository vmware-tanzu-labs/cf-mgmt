&larr; [back to Commands](../README.md)

# `cf-mgmt-config generate-concourse-pipeline`

`generate-concourse-pipeline` generate a pipeline.yml, vars.yml and necessary task yml files for running all the tasks listed below.  Just need to update your vars.yml and check in all your code to GIT and execute the fly command to register your pipeline. ```vars.yml``` contains place holders for LDAP and CF user credentials. If you do not prefer storing the credentials in ```vars.yml```, you can pass them via the ```fly``` command line arguments.

## Command Usage

```
Usage:
  cf-mgmt-config [OPTIONS] generate-concourse-pipeline [generate-concourse-pipeline-OPTIONS]

Help Options:
  -h, --help            Show this help message

[generate-concourse-pipeline command options]
  --config-dir= Name of the config directory (default: config) [$CONFIG_DIR]
  --target-dir= Name of the target directory to generate into (default: .)
```

Once the pipeline files are generated, you can create a pipeline as follows:

```
fly -t  login -c <concourse_instance>
fly -t <targetname> set-pipeline -p <pipeline_name> \
   -c pipeline.yml \
   -l vars.yml \
   —-var "ldap_password=<ldap_password>" \
   --var "client_secret=<client_secret>" \
   —-var "password=<org/space_admin_password>"
```
Note: If using a _client_secret_ set your _password_ in the ```vars.yml``` to an empty string ```""```

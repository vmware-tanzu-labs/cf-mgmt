&larr; [back to Commands](../README.md)

# `cf-mgmt-config generate-concourse-pipeline`

`generate-concourse-pipeline` generate a pipeline.yml, vars.yml and necessary task yml files for running all the tasks listed below.  Just need to update your vars.yml and check in all your code to GIT and execute the fly command to register your pipeline. ```vars.yml``` contains place holders for LDAP and CF user credentials. If you do not prefer storing the credentials in ```vars.yml```, you can pass them via the ```fly``` command line arguments.

## Command Usage

```
Usage:
  main [OPTIONS] generate-concourse-pipeline

Help Options:
  -h, --help      Show this help message
```

Once the pipeline files are generated, you can create a pipeline as follows:

```
fly -t  login -c <concourse_instance>
fly -t <targetname> set-pipeline -p <pipeline_name> \
   -c pipeline.yml \
   -l vars.yml \
   —-var "ldap_password=<ldap_password>" \
   --var "client_secret=<client_sercret>" \
   —-var "password=<org/space_admin_password>"
```

## Using the data source
Example exists in `alert-configurations-data.tf`. To use this example exactly:
- Copy directory to local disk
- Add a `terraform.tfvars`
- Add your `project_id`
- Run `terraform apply`

### Create alert resources and import them into state file
```
terraform output -raw alert_imports > import-alerts.sh
terraform output -raw alert_resources > alert-configurations.tf
chmod +x ./import-alerts.sh
./import-alerts.sh
terraform apply
```

## Contingency Plans
If unhappy with the resource file or imports, here are some things that can be done:

### Remove targeted resources from the appropriate files and remove the alet_configuration from state
- Manually remove the resource (ex: `mongodbatlas_alert_configuration.CLUSTER_MONGOS_IS_MISSING_2`) from the `tf` file, and then remove it from state, ex:
```
terraform state rm mongodbatlas_alert_configuration.CLUSTER_MONGOS_IS_MISSING_2
```

### Remove all alert_configurations from state
- Delete the `tf` file that was used for import, and then:
```
terraform state list | grep ^mongodbatlas_alert_configuration. | awk '{print "terraform state rm " $1}' > state-rm-alerts.sh
chmod +x state-rm-alerts.sh
./state-rm-alerts.sh
```

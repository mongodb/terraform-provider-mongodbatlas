# MongoDB Atlas Cluster Adaptive Settings

This example manages Adaptive Settings overrides for an existing MongoDB Atlas cluster and reads the resulting effective settings with the corresponding data source.

The example enables both overrides by default. Set `load_shedding_enabled` or `search_load_shedding_enabled` to `false` to disable the corresponding behavior. The API keys are `LOAD_SHEDDING` and `SEARCH_OVERLOAD_PROTECTION`; Search Load Shedding retains the latter key.

Configure MongoDB Atlas provider credentials with environment variables, then provide the project ID and cluster name:

```hcl
project_id   = "0123456789abcdef01234567"
cluster_name = "Cluster0"
```

Run the example with:

```bash
terraform init
terraform apply
```

Destroy the configuration when you finish so Atlas resets the overrides to its managed defaults:

```bash
terraform destroy
```

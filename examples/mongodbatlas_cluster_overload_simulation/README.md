# MongoDB Atlas Cluster Overload Simulation

This example creates a cluster, starts an overload protection simulation, and reads the simulation with the corresponding data source.

Configure MongoDB Atlas provider credentials with environment variables, then supply the project ID:

```hcl
project_id = "0123456789abcdef01234567"
```

Run the example with:

```bash
terraform init
terraform apply
```

Destroy the configuration when you finish testing so the simulation is cancelled and the cluster is removed:

```bash
terraform destroy
```

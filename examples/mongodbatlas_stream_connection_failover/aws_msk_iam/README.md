# AWS MSK IAM failover connection

Create a stream workspace with a failover region and an Amazon MSK connection that authenticates with `AWS_MSK_IAM` in both the default and failover regions.

Both connections set `authentication.aws.role_arn`. The role must already be registered through [Atlas Cloud Provider Access](https://www.mongodb.com/docs/atlas/security/set-up-unified-aws-access/) so that Atlas can assume it.

## Prerequisites

- Project Owner on the project.
- An Amazon MSK cluster in each region, with IAM access control enabled.
- An AWS IAM role trusted by Atlas Cloud Provider Access.

Credentials come from the environment (`MONGODB_ATLAS_CLIENT_ID` and `MONGODB_ATLAS_CLIENT_SECRET`).

## Usage

**1. Set credentials.**

```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

**2. Create `terraform.tfvars`.**

```hcl
project_id                 = "your-project-id"
bootstrap_servers          = "b-1.example.kafka.us-east-1.amazonaws.com:9098,b-2.example.kafka.us-east-1.amazonaws.com:9098"
failover_bootstrap_servers = "b-1.example.kafka.eu-west-1.amazonaws.com:9098"
kafka_iam_role_arn         = "arn:aws:iam::123456789012:role/atlas-msk-access"
```

The workspace targets `VIRGINIA_USA` as its default region and `DUBLIN_IRL` as its failover region. Change both in `main.tf` to match where your MSK clusters run.

**3. Plan and apply.**

```bash
terraform plan
terraform apply
```

**4. Destroy.**

```bash
terraform destroy
```

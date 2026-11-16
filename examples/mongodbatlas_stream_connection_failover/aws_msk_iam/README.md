# MongoDB Atlas Provider — Stream Connection Failover with AWS MSK IAM authentication

This example creates a stream workspace with a failover region, a primary `mongodbatlas_stream_connection`
authenticated to an Amazon MSK cluster with `AWS_MSK_IAM`, and a `mongodbatlas_stream_connection_failover`
for the failover region that authenticates to the cluster in that region with the same IAM role.

Both connections set `authentication.aws.role_arn`. The role must already be registered through
[Atlas Cloud Provider Access](https://www.mongodb.com/docs/atlas/security/set-up-unified-aws-access/)
so that Atlas can assume it.

## Dependencies

* Terraform MongoDB Atlas Provider
* A MongoDB Atlas account
* An Amazon MSK cluster in each region, with IAM access control enabled
* An AWS IAM role trusted by Atlas Cloud Provider Access

## Usage

```bash
terraform init
terraform plan -var="project_id=<your-project-id>" \
  -var="bootstrap_servers=<primary-region-brokers>" \
  -var="failover_bootstrap_servers=<failover-region-brokers>" \
  -var="kafka_iam_role_arn=<role-arn>"
terraform apply
```

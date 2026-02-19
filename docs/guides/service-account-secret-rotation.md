---
page_title: "Guide: Service Account Secret Rotation"
---

# Guide: Service Account Secret Rotation

**Objective**: This guide shows a simple approach to manage and rotate Service Account secrets via Terraform, using a two-secret rotation pattern that allows you to rotate secrets without downtime.

## Overview

When you create a Service Account, Atlas automatically generates a secret. The secret value is returned only once, at creation time.

For production environments, you typically want to maintain two secrets at any given time, allowing you to rotate one while the other remains active.

This guide applies to both organization-level and project-level service accounts:
- **Organization-level**: Use `mongodbatlas_service_account` and `mongodbatlas_service_account_secret`
- **Project-level**: Use `mongodbatlas_project_service_account` and `mongodbatlas_project_service_account_secret`

**Note**: The steps below use organization-level resources, but the same approach applies to project-level resources.

~> **WARNING:** Service Account secrets expire after the configured `secret_expires_after_hours` period. To avoid losing access to the Atlas Administration API, update your application with the new client secret as soon as possible after rotation. If all secrets expire before being replaced, you will lose access to the organization. For more information, see [Rotate Service Account Secrets](https://www.mongodb.com/docs/atlas/tutorial/rotate-service-account-secrets/).

## Best Practices Before Starting

- **Backup your Terraform state file** before making any changes.
- **Test the rotation process in a non-production environment** if possible.
- Managing Service Accounts with Terraform **exposes sensitive organizational secrets** in Terraform's state. We suggest following [Terraform's best practices](https://developer.hashicorp.com/terraform/language/state/sensitive-data).

## Setup

### Step 1: Initial Configuration

1. Start with the following configuration. It creates a service account (which includes an initial secret) and a second secret:

```terraform
variable "org_id" {
  description = "MongoDB Atlas Organization ID"
  type        = string
}

# Create service account (also creates the first secret)
resource "mongodbatlas_service_account" "this" {
  org_id                     = var.org_id
  name                       = "example-service-account"
  description                = "Example Service Account"
  roles                      = ["ORG_READ_ONLY"]
  secret_expires_after_hours = 2160 # 90 days
}

# Create secret_2 as a separate resource
resource "mongodbatlas_service_account_secret" "secret_2" {
  org_id                     = var.org_id
  client_id                  = mongodbatlas_service_account.this.client_id
  secret_expires_after_hours = 2160 # 90 days
}

# Output the import ID for secret_1
output "secret_1_import_id" {
  value       = "${var.org_id}/${mongodbatlas_service_account.this.client_id}/${mongodbatlas_service_account.this.secrets[0].secret_id}"
  description = "Import ID for secret_1. Use this to import the initial secret into Terraform."
}

output "secret_2" {
  sensitive = true
  value     = mongodbatlas_service_account_secret.secret_2.secret
}
```

2. Apply the configuration:
```shell
terraform apply
```

3. Copy the `secret_1_import_id` value from the output. It is required for Step 2.

4. Retrieve and securely store the `secret_2` value (**warning**: this prints the secret to your terminal):
```shell
terraform output -raw secret_2
```

### Step 2: Import the Initial Secret into Terraform

To manage the initial secret (created automatically with the service account) as a Terraform resource, you need to import it.

1. Add the `secret_1` resource and output to your configuration:
```terraform
# Define secret_1
resource "mongodbatlas_service_account_secret" "secret_1" {
  org_id    = var.org_id
  client_id = mongodbatlas_service_account.this.client_id
}

output "secret_1" {
  sensitive = true
  value     = mongodbatlas_service_account_secret.secret_1.secret
}
```

2. Import `secret_1` using the ID from the previous step:

```shell
terraform import mongodbatlas_service_account_secret.secret_1 <secret_1_import_id>
```

**Note**: After import, `mongodbatlas_service_account_secret.secret_1.secret` is `null` since secret values are only returned at creation time. The secret will have a value after the first rotation.

3. Verify that the import was successful:
```shell
terraform plan
```

You should see no planned changes.

4. Remove the `secret_1_import_id` output. It is no longer needed.

## Secret Rotation

After the initial setup is complete, you can rotate secrets using Terraform's `-replace` flag. This recreates the resource, generating a new secret.

### Rotate secret_1

1. Add the `secret_expires_after_hours` attribute to the `secret_1` resource:
```terraform
resource "mongodbatlas_service_account_secret" "secret_1" {
  org_id                     = var.org_id
  client_id                  = mongodbatlas_service_account.this.client_id
  secret_expires_after_hours = 2160 # 90 days
}
```

2. Rotate the secret:
```shell
terraform apply -replace="mongodbatlas_service_account_secret.secret_1"
```

3. Retrieve and securely store the new secret value (**warning**: this prints the secret to your terminal):
```shell
terraform output -raw secret_1
```

4. Update your applications with the new secret value.

### Rotate secret_2

**Note**: `secret_2` already has a value from the initial setup. You can skip this section until you need to rotate it.

1. Rotate the secret:

```shell
terraform apply -replace="mongodbatlas_service_account_secret.secret_2"
```

2. Retrieve and securely store the new secret value (**warning**: this prints the secret to your terminal):

```shell
terraform output -raw secret_2
```

3. Update your applications with the new secret value.

### Ongoing rotation

Continue alternating between secrets when rotating. This ensures that the older secret remains active while the new one is rotated and deployed to your applications.

## Complete Configuration

This is the full configuration after the first rotation is complete:

```terraform
variable "org_id" {
  description = "MongoDB Atlas Organization ID"
  type        = string
}

resource "mongodbatlas_service_account" "this" {
  org_id                     = var.org_id
  name                       = "example-service-account"
  description                = "Example Service Account"
  roles                      = ["ORG_READ_ONLY"]
  secret_expires_after_hours = 2160 # 90 days
}

resource "mongodbatlas_service_account_secret" "secret_1" {
  org_id                     = var.org_id
  client_id                  = mongodbatlas_service_account.this.client_id
  secret_expires_after_hours = 2160 # 90 days
}

resource "mongodbatlas_service_account_secret" "secret_2" {
  org_id                     = var.org_id
  client_id                  = mongodbatlas_service_account.this.client_id
  secret_expires_after_hours = 2160 # 90 days
}

output "secret_1" {
  sensitive = true
  value     = mongodbatlas_service_account_secret.secret_1.secret
}

output "secret_2" {
  sensitive = true
  value     = mongodbatlas_service_account_secret.secret_2.secret
}
```

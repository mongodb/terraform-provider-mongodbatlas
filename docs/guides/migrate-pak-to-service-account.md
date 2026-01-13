---
page_title: "Migration Guide: Programmatic API Keys (PAKs) to Service Accounts (SAs)"
---

# Migration Guide: Programmatic API Keys (PAKs) to Service Accounts (SAs)

## Overview

This guide explains how to migrate from Programmatic API Key (PAK) resources to Service Account (SA) resources.

**Note:** Migration to Service Accounts is recommended but **not required**. If you are currently using API Key resources, you may continue to do so. This guide is for users who wish to adopt Service Accounts for greater security or best practices, but existing PAK configurations will continue to work and be supported.

## Before You Begin

- **Backup your Terraform state file** before making any changes.
- **Test the process in a non-production environment** if possible.
- Managing Service Accounts with Terraform **exposes sensitive organizational secrets** in Terraform's state. We suggest following [Terraform's best practices](https://developer.hashicorp.com/terraform/language/state/sensitive-data).


---

<details>
  <summary><span style="font-size:1.4em; font-weight:bold;">Organization-Level Migration</span></summary>

## Organization-Level API Keys to Service Accounts

**Objective**: Migrate from organization-level PAK resources (`mongodbatlas_api_key`, `mongodbatlas_api_key_project_assignment`, `mongodbatlas_access_list_api_key`) to organization-level Service Account resources (`mongodbatlas_service_account`, `mongodbatlas_service_account_project_assignment`, `mongodbatlas_service_account_access_list_entry`).

### Resource Mapping
The following table shows the mapping between organization-level PAK resources and their Service Account equivalents:

| PAK Resource | Service Account Resource | Notes |
|--------------|-------------------------|-------|
| `mongodbatlas_api_key` | `mongodbatlas_service_account` | API key / Service Account |
| `mongodbatlas_api_key_project_assignment` | `mongodbatlas_service_account_project_assignment` | Project assignment |
| `mongodbatlas_access_list_api_key` | `mongodbatlas_service_account_access_list_entry` | IP access list entry |

---

### Migration Steps

For complete working examples, see the [organization-level migration example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/migrate_pak_to_service_account/org_level).

### Step 1: Initial Configuration - PAK Resources Only

Original configuration with PAK resources:

```terraform
resource "mongodbatlas_api_key" "this" {
  org_id      = var.org_id
  description = "Example API Key"
  role_names  = ["ORG_MEMBER"]
}

# Project assignment for the API Key
resource "mongodbatlas_api_key_project_assignment" "this" {
  project_id = var.project_id
  api_key_id = mongodbatlas_api_key.this.api_key_id
  roles      = ["GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_ONLY"]
}

# IP Access List entry for the API Key
resource "mongodbatlas_access_list_api_key" "this" {
  org_id     = var.org_id
  api_key_id = mongodbatlas_api_key.this.api_key_id
  cidr_block = "192.168.1.100/32"
  # Alternative: ip_address = "192.168.1.100"
}
```

### Step 2: Intermediate State - Add Service Account Resources Alongside Existing PAK Resources

Add the Service Account resources to your configuration while keeping the existing PAK resources. This allows both authentication methods to work simultaneously, enabling you to test Service Accounts before removing PAKs.

```terraform
resource "mongodbatlas_service_account" "this" {
  org_id                     = var.org_id
  name                       = "example-service-account"
  description                = "Example Service Account"
  roles                      = ["ORG_MEMBER"]
  secret_expires_after_hours = 2160 # 90 days
}

resource "mongodbatlas_service_account_project_assignment" "this" {
  project_id = var.project_id
  client_id  = mongodbatlas_service_account.this.client_id
  roles      = ["GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_ONLY"]
}

resource "mongodbatlas_service_account_access_list_entry" "this" {
  org_id     = var.org_id
  client_id  = mongodbatlas_service_account.this.client_id
  cidr_block = "192.168.1.100/32"
  # Alternative: ip_address = "192.168.1.100"
}

output "service_account_first_secret" {
  description = "The secret value of the first secret created with the service account. Only available after initial creation."
  value       = try(mongodbatlas_service_account.this.secrets[0].secret, null)
  sensitive   = true
}
```

**Apply and test:**

1. Run `terraform plan` to review the changes.
2. Run `terraform apply` to create the Service Account resources.
3. Retrieve and securely store the `service_account_first_secret` value (**warning**: this prints the secret to your terminal):

   ```bash
   terraform output -raw service_account_first_secret
   ```

4. Test your Service Account in your applications and verify that both PAK and SA authentication methods work correctly.
5. Re-run `terraform plan` to ensure you have no unexpected changes: `No changes. Your infrastructure matches the configuration.`

### Step 3: Final State - Remove PAK Resources, SA Resources Only

Once you have verified that the Service Account works correctly, remove the PAK resources from your configuration:

```terraform
resource "mongodbatlas_service_account" "this" {
  org_id                     = var.org_id
  name                       = "example-service-account"
  description                = "Example Service Account"
  roles                      = ["ORG_MEMBER"]
  secret_expires_after_hours = 2160 # 90 days
}

resource "mongodbatlas_service_account_project_assignment" "this" {
  project_id = var.project_id
  client_id  = mongodbatlas_service_account.this.client_id
  roles      = ["GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_ONLY"]
}

resource "mongodbatlas_service_account_access_list_entry" "this" {
  org_id     = var.org_id
  client_id  = mongodbatlas_service_account.this.client_id
  cidr_block = "192.168.1.100/32"
}
```

1. Run `terraform plan` to verify:
   - PAK resources are planned for destruction
   - Only the Service Account resources remain
   - No unexpected changes

2. Run `terraform apply` to finalize the migration. This will delete the PAK resources from Atlas.

3. Verify that your applications and infrastructure continue to work with Service Accounts.

4. Re-run `terraform plan` to ensure you have no planned changes: `No changes. Your infrastructure matches the configuration.`


</details>

---

<details>
  <summary><span style="font-size:1.4em; font-weight:bold;">Project-Level Migration</span></summary>

## Project-Level API Keys to Service Accounts

**Objective**: Migrate from project-level PAK resources (`mongodbatlas_project_api_key`, `mongodbatlas_access_list_api_key`) to project-level Service Account resources (`mongodbatlas_project_service_account`, `mongodbatlas_project_service_account_access_list_entry`).

**Important:** Organization-level resources (`mongodbatlas_service_account`) are the recommended approach. Project-level resources (`mongodbatlas_project_service_account`) should only be used if you do not have organization-level permissions to manage service accounts. Otherwise, use the [Organization-Level Migration](#organization-level-api-keys-to-service-accounts) approach.

### Resource Mapping

The following table shows the mapping between project-level PAK resources and their Service Account equivalents:

| PAK Resource | Service Account Resource | Notes |
|--------------|-------------------------|-------|
| `mongodbatlas_project_api_key` | `mongodbatlas_project_service_account` | Project-level API key / Service Account (use only if you don't have org-level permissions) |
| `mongodbatlas_access_list_api_key` | `mongodbatlas_project_service_account_access_list_entry` | IP access list entry for project-level service account |

**Important:** Organization-level resources (`mongodbatlas_service_account`) are the recommended approach. Project-level resources (`mongodbatlas_project_service_account`) should only be used if you do not have organization-level permissions to manage service accounts.

---

### Migration Steps

For complete working examples, see the [project-level migration example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/migrate_pak_to_service_account/project_level).

### Step 1: Initial Configuration - PAK Resources Only

Original configuration with PAK resources:

```terraform
resource "mongodbatlas_project_api_key" "this" {
  description = "Example Project API Key"
  project_assignment {
    project_id = var.project_id
    role_names = ["GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_ONLY"]
  }
}

# IP Access List entry for the API Key
resource "mongodbatlas_access_list_api_key" "this" {
  org_id     = var.org_id
  api_key_id = mongodbatlas_project_api_key.this.api_key_id
  cidr_block = "192.168.1.100/32"
  # Alternative: ip_address = "192.168.1.100"
}
```

### Step 2: Intermediate State - Add Service Account Resources Alongside Existing PAK Resources

Add the Service Account resources to your configuration while keeping the existing PAK resources. This allows both authentication methods to work simultaneously, enabling you to test Service Accounts before removing PAKs.

```terraform
resource "mongodbatlas_project_service_account" "this" {
  project_id                 = var.project_id
  name                       = "example-project-service-account"
  description                = "Example Project Service Account"
  roles                      = ["GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_ONLY"]
  secret_expires_after_hours = 2160 # 90 days
}

resource "mongodbatlas_project_service_account_access_list_entry" "this" {
  project_id = var.project_id
  client_id  = mongodbatlas_project_service_account.this.client_id
  cidr_block = "192.168.1.100/32"
  # Alternative: ip_address = "192.168.1.100"
}

output "project_service_account_first_secret" {
  description = "The secret value of the first secret created with the project service account. Only available after initial creation."
  value       = try(mongodbatlas_project_service_account.this.secrets[0].secret, null)
  sensitive   = true
}
```

**Apply and test:**

1. Run `terraform plan` to review the changes.
2. Run `terraform apply` to create the Service Account resource.
3. Retrieve and securely store the `project_service_account_first_secret` value (**warning**: this prints the secret to your terminal):

   ```bash
   terraform output -raw project_service_account_first_secret
   ```

4. Test your Service Account in your applications and verify that both PAK and SA authentication methods work correctly.
5. Re-run `terraform plan` to ensure you have no unexpected changes: `No changes. Your infrastructure matches the configuration.`

### Step 3: Final State - Remove PAK Resources, SA Resources Only

Once you have verified that the Service Account works correctly, remove the PAK resources from your configuration:

```terraform
resource "mongodbatlas_project_service_account" "this" {
  project_id                 = var.project_id
  name                       = "example-project-service-account"
  description                = "Example Project Service Account"
  roles                      = ["GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_ONLY"]
  secret_expires_after_hours = 2160 # 90 days
}

resource "mongodbatlas_project_service_account_access_list_entry" "this" {
  project_id = var.project_id
  client_id  = mongodbatlas_project_service_account.this.client_id
  cidr_block = "192.168.1.100/32"
}
```

1. Run `terraform plan` to verify:
   - PAK resources are planned for destruction
   - Only the Service Account resources remain
   - No unexpected changes

2. Run `terraform apply` to finalize the migration. This will delete the PAK resource from Atlas.

3. Verify that your applications and infrastructure continue to work with Service Accounts.

4. Re-run `terraform plan` to ensure you have no planned changes: `No changes. Your infrastructure matches the configuration.`

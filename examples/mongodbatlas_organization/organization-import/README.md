# MongoDB Atlas Organization Import Example

This example demonstrates how to import an existing MongoDB Atlas organization into Terraform management.

## Overview

When you have an existing MongoDB Atlas organization that you want to manage with Terraform, you can import it using the organization ID. This is useful when you have an existing organization that was created outside of Terraform and want to start managing its settings with Infrastructure as Code.

Don't include `org_owner_id`, `description`, `role_names` or `federation_settings_id` when importing an organization as they are creation-only attributes.

## Dependencies

- Terraform v0.15 or greater
- provider.mongodbatlas: version = "~> 1.38.0"
- MongoDB Atlas account with an existing organization
- API key with Organization Owner permissions for the organization you want to import
- Organization ID of the existing organization

## Usage

### 1. Set up your credentials

Create a `terraform.tfvars` file with your Atlas credentials:

```hcl
public_key  = "<PUBLIC_KEY>"
private_key = "<PRIVATE_KEY>"
org_name    = "<ORG_NAME>"

# Optional: Configure organization settings
api_access_list_required   = false
multi_factor_auth_required = true
restrict_employee_access   = true
gen_ai_features_enabled    = true
security_contact          = "<SECURITY_CONTACT_EMAIL>"
skip_default_alerts_settings = true
```

### 2. Initialize Terraform

```bash
terraform init
```

### 3. Import the organization

Replace `<ORG_ID>` with your actual organization ID:

```bash
terraform import mongodbatlas_organization.imported <ORG_ID>
```

To find your organization ID:
1. Log into MongoDB Atlas
2. Go to your organization settings
3. The organization ID is displayed in the URL or settings page

### 4. Review the plan

```bash
terraform plan
```

This will show you what changes (if any) Terraform will make to align your organization with the configuration.

### 5. Apply changes

```bash
terraform apply
```

## Finding Your Organization ID

You can find your organization ID in several ways:

1. **Atlas UI**: Navigate to your organization settings - the ID is visible in the URL
2. **Atlas API**: Use the `organizations` endpoint to list all organizations
3. **Atlas CLI**: Use `atlas organizations list` if you have the Atlas CLI installed

## Outputs

After successful import, you'll have access to:
- `org_name` - The organization name

**Note**: Unlike when creating a new organization, importing does not provide `public_key` and `private_key` outputs since these are only generated during organization creation.

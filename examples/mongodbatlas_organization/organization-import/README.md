# MongoDB Atlas Organization Import Example

This example demonstrates how to import an existing MongoDB Atlas organization into Terraform management.

## Overview

When you have an existing MongoDB Atlas organization that you want to manage with Terraform, you can import it using the organization ID. This is useful when:

- You have an existing organization that was created outside of Terraform
- You want to start managing an organization's settings with Infrastructure as Code
- You need to bring existing resources under Terraform management

## Important Notes

### Creation-Only Attributes

**DO NOT** include these attributes when importing an organization:
- `org_owner_id` - Only used during creation
- `description` - Only used during creation
- `role_names` - Only used during creation

These attributes are creation-only and will cause errors if specified during import.

### Import-Compatible Attributes

You can configure these attributes when importing:
- `name` - Required (must match existing organization name)
- `api_access_list_required` - Optional
- `multi_factor_auth_required` - Optional
- `restrict_employee_access` - Optional
- `gen_ai_features_enabled` - Optional
- `security_contact` - Optional
- `skip_default_alerts_settings` - Optional
- `federation_settings_id` - Optional (if using federation)

## Prerequisites

- Terraform v0.13 or greater
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

## Important Warnings

### Destroy Behavior

When you run `terraform destroy` on an imported organization:
- The organization will be **removed from Terraform state**
- The actual organization in MongoDB Atlas will **NOT be deleted**
- This is a safety measure to prevent accidental deletion of production organizations

### State Management

- The imported organization will be managed by Terraform going forward
- Any manual changes made outside of Terraform may cause drift
- Always review plans before applying to avoid unexpected changes

## Troubleshooting

### Common Issues

1. **"org_owner_id cannot be changed after creation"**
   - Remove `org_owner_id` from your configuration - it's creation-only

2. **"description cannot be changed after creation"**
   - Remove `description` from your configuration - it's creation-only

3. **"role_names cannot be changed after creation"**
   - Remove `role_names` from your configuration - it's creation-only

4. **Import fails with "not found"**
   - Verify the organization ID is correct
   - Ensure your API key has access to the organization

### Getting Help

If you encounter issues:
1. Check the [Terraform Provider Documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs)
2. Review the [MongoDB Atlas API Documentation](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/)
3. Open an issue on the [Provider GitHub Repository](https://github.com/mongodb/terraform-provider-mongodbatlas) 

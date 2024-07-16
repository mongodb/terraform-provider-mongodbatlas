# Data Source: mongodbatlas_project_api_key

`mongodbatlas_project_api_key` describes a MongoDB Atlas Project API Key. This represents a Project API Key that has been created.

~> **IMPORTANT WARNING:** Managing Atlas Programmatic API Keys (PAKs) with Terraform will expose sensitive organizational secrets in Terraform's state. We suggest following [Terraform's best practices](https://developer.hashicorp.com/terraform/language/state/sensitive-data). You may also want to consider managing your PAKs via a more secure method, such as the [HashiCorp Vault MongoDB Atlas Secrets Engine](https://developer.hashicorp.com/vault/docs/secrets/mongodbatlas).

-> **NOTE:** You may find project_id in the official documentation.

## Example Usage

### Using project_id and api_key_id attribute to query
```terraform

resource "mongodbatlas_project_api_key" "test" {
  description   = "Description of your API key"
  project_assignment {
    project_id = "64259ee860c43338194b0f8e"
    role_names = ["GROUP_READ_ONLY"]
  }
}

data "mongodbatlas_project_api_key" "test" {
  project_id = "64259ee860c43338194b0f8e"
  api_key_id = mongodbatlas_api_key.test.api_key_id
}
```

## Argument Reference

* `project_id` - (Required) The unique ID for the project.
* `api_key_id` - (Required) Unique identifier for this Project API key.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `description` - Description of this Project API key.
* `public_key` - Public key for this Organization API key.
* `private_key` - Private key for this Organization API key.

### project_assignment
List of Project roles that the Programmatic API key needs to have.

* `project_id` -  Project ID to assign to Access Key
* `role_names` -  List of Project roles that the Programmatic API key needs to have. Ensure you provide: at least one role and ensure all roles are valid for the Project. You must specify an array even if you are only associating a single role with the Programmatic API key. The [MongoDB Documentation](https://www.mongodb.com/docs/atlas/reference/user-roles/#project-roles) describes the valid roles that can be assigned.
    

See [MongoDB Atlas API - API Key](https://www.mongodb.com/docs/atlas/reference/api/projectApiKeys/get-all-apiKeys-in-one-project/) - Documentation for more information.

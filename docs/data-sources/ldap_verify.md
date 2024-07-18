# Data Source: mongodbatlas_ldap_verify

`mongodbatlas_ldap_verify` describes a LDAP Verify.

-> **NOTE:** Groups and projects are synonymous terms. You may find **group_id** in the official documentation.


## Example Usage

```terraform
resource "mongodbatlas_project" "test" {
  name   = "NAME OF THE PROJECT"
  org_id = "ORG ID"
}

resource "mongodbatlas_advanced_cluster" "test" {
  project_id     = mongodbatlas_project.test.id
  name           = "ClusterName"
  cluster_type   = "REPLICASET"
  backup_enabled = true # enable cloud provider snapshots

  replication_specs {
    region_configs {
      priority      = 7
      provider_name = "AWS"
      region_name   = "US_EAST_1"
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
    }
  }
}

resource "mongodbatlas_ldap_verify" "test" {
  project_id                  = mongodbatlas_project.test.id
  hostname = "HOSTNAME"
  port                     = 636
  bind_username                     = "USERNAME"
  bind_password                     = "PASSWORD"
  depends_on = [mongodbatlas_advanced_cluster.test]
}

data "mongodbatlas_ldap_verify" "test" {
  project_id = mongodbatlas_ldap_verify.test.project_id
  request_id = mongodbatlas_ldap_verify.test.request_id
}
```

## Argument Reference

* `project_id` - (Required) Unique identifier for the Atlas project associated with the verification request.
* `request_id` - (Required) Unique identifier of a request to verify an LDAP configuration.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `hostname` - (Required) The hostname or IP address of the LDAP server.
* `port` - LDAP ConfigurationThe port to which the LDAP server listens for client connections.
* `bind_username` - The user DN that Atlas uses to connect to the LDAP server.
* `bind_password` - The password used to authenticate the `bind_username`.
* `ca_certificate` - LDAP ConfigurationCA certificate used to verify the identify of the LDAP server.
* `authz_query_template` - LDAP ConfigurationAn LDAP query template that Atlas executes to obtain the LDAP groups to which the authenticated user belongs.
* `request_id` - The unique identifier for the request to verify the LDAP over TLS/SSL configuration.
* `status` - The current status of the LDAP over TLS/SSL configuration.
* `links` - One or more links to sub-resources. The relations in the URLs are explained in the Web Linking Specification.
* `validations` - Array of validation messages related to the verification of the provided LDAP over TLS/SSL configuration details.


See detailed information for arguments and attributes: [MongoDB API LDAP Verify](https://docs.atlas.mongodb.com/reference/api/ldaps-configuration-verification-status)
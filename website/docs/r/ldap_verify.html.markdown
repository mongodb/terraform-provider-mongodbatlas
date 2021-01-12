---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: ldap-verify"
sidebar_current: "docs-mongodbatlas-resource-ldap-verify"
description: |-
    Provides a LDAP Verify resource.
---

# mongodbatlas_ldap_verify

`mongodbatlas_ldap_verify` provides an LDAP Verify resource. This allows ldap verify to be created.

## Example Usage

```hcl
resource "mongodbatlas_project" "test" {
	name   = "NAME OF THE PROJECT"
	org_id = "ORG ID"
}

resource "mongodbatlas_cluster" "test" {
    project_id   = mongodbatlas_project.test.id
    name         = "NAME OF THE CLUSTER"
    disk_size_gb = 5
    
    // Provider Settings "block"
    provider_name               = "AWS"
    provider_region_name        = "US_EAST_2"
    provider_instance_size_name = "M10"
    provider_backup_enabled     = true //enable cloud provider snapshots
    provider_disk_iops          = 100
    provider_encrypt_ebs_volume = false
}

resource "mongodbatlas_ldap_verify" "test" {
    project_id                  = mongodbatlas_project.test.id
    hostname = "HOSTNAME"
    port                     = 636
    bind_username                     = "USERNAME"
    bind_password                     = "PASSWORD"
    depends_on = ["mongodbatlas_cluster.test"]
}
```

## Argument Reference

* `project_id` - (Required) The unique ID for the project to configure LDAP.
* `hostname` - (Required) The hostname or IP address of the LDAP server. The server must be visible to the internet or connected to your Atlas cluster with VPC Peering.
* `port` - (Optional) The port to which the LDAP server listens for client connections. Default: `636`
* `bind_username` - (Required) The user DN that Atlas uses to connect to the LDAP server. Must be the full DN, such as `CN=BindUser,CN=Users,DC=myldapserver,DC=mycompany,DC=com`.
* `bind_password` - (Required) The password used to authenticate the `bind_username`.
* `ca_certificate` - (Optional) CA certificate used to verify the identify of the LDAP server. Self-signed certificates are allowed.
* `authz_query_template` - (Optional) An LDAP query template that Atlas executes to obtain the LDAP groups to which the authenticated user belongs. Used only for user authorization. Use the {USER} placeholder in the URL to substitute the authenticated username. The query is relative to the host specified with hostname. The formatting for the query must conform to RFC4515 and RFC 4516. If you do not provide a query template, Atlas attempts to use the default value: `{USER}?memberOf?base`.

~> **NOTE:** LDAP Configuration created by API Keys must belong to an existing organization.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `request_id` - The unique identifier for the request to verify the LDAP over TLS/SSL configuration.
* `status` - The current status of the LDAP over TLS/SSL configuration. One of the following values: `PENDING`, `SUCCESS`, and `FAILED`.
* `links` - One or more links to sub-resources. The relations in the URLs are explained in the Web Linking Specification.
* `validations` - Array of validation messages related to the verification of the provided LDAP over TLS/SSL configuration details. The array contains a document for each test that Atlas runs. Atlas stops running tests after the first failure. The following return values can be seen here: [Values](https://docs.atlas.mongodb.com/reference/api/ldaps-configuration-request-verification)
    
## Import

LDAP Configuration must be imported using project ID and request ID, e.g.

```
$ terraform import mongodbatlas_ldap_verify.test 5d09d6a59ccf6445652a444a-5d09d6a59ccf6445652a444a
```

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/ldaps-configuration-request-verification)
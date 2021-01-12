---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: ldap configuration"
sidebar_current: "docs-mongodbatlas-datasource-ldap-configuration"
description: |-
    Describes a LDAP Configuration.
---

# mongodbatlas_ldap_configuration

`mongodbatlas_ldap_configuration` describes a LDAP Configuration.

-> **NOTE:** Groups and projects are synonymous terms. You may find **group_id** in the official documentation.


## Example Usage

```hcl
resource "mongodbatlas_project" "test" {
  name   = "NAME OF THE PROJECT"
  org_id = "ORG ID"
}

resource "mongodbatlas_ldap_configuration" "test" {
  project_id                  = mongodbatlas_project.test.id
  authentication_enabled      = true
  hostname 					= "HOSTNAME"
  port                     	= 636
  bind_username               = "USERNAME"
  bind_password               = "PASSWORD"
}

data "mongodbatlas_ldap_configuration" "test" {
  project_id = mongodbatlas_ldap_configuration.test.id
}
```

## Argument Reference

* `project_id` - (Required) Identifier for the Atlas project associated with the LDAP over TLS/SSL configuration.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `authentication_enabled` - Specifies whether user authentication with LDAP is enabled.
* `authorization_enabled` - Specifies whether user authorization with LDAP is enabled.
* `hostname` - (Required) The hostname or IP address of the LDAP server.
* `port` - LDAP ConfigurationThe port to which the LDAP server listens for client connections.
* `bind_username` - The user DN that Atlas uses to connect to the LDAP server.
* `bind_password` - The password used to authenticate the `bind_username`.
* `ca_certificate` - LDAP ConfigurationCA certificate used to verify the identify of the LDAP server.
* `authz_query_template` - LDAP ConfigurationAn LDAP query template that Atlas executes to obtain the LDAP groups to which the authenticated user belongs.
* `user_to_dn_mapping` - LDAP ConfigurationMaps an LDAP username for authentication to an LDAP Distinguished Name (DN).
* `user_to_dn_mapping.0.match` - LDAP ConfigurationA regular expression to match against a provided LDAP username.
* `user_to_dn_mapping.0.substitution` - LDAP ConfigurationAn LDAP Distinguished Name (DN) formatting template that converts the LDAP name matched by the `match` regular expression into an LDAP Distinguished Name.
* `user_to_dn_mapping.0.ldap_query` - LDAP ConfigurationAn LDAP query formatting template that inserts the LDAP name matched by the `match` regular expression into an LDAP query URI as specified by RFC 4515 and RFC 4516.


See detailed information for arguments and attributes: [MongoDB API LDAP Configuration](https://docs.atlas.mongodb.com/reference/api/ldaps-configuration-get-current)
---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: ldap-configuration"
sidebar_current: "docs-mongodbatlas-resource-ldap-configuration"
description: |-
    Provides a LDAP Configuration resource.
---

# Resource: mongodbatlas_ldap_configuration

`mongodbatlas_ldap_configuration` provides an LDAP Configuration resource. This allows an LDAP configuration for an Atlas project to be crated and managed. This endpoint doesn’t verify connectivity using the provided LDAP over TLS configuration details. To verify a configuration before saving it, use the resource to [verify](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/INTMDB-114/website/docs/r/ldap_verify.html.markdown) the LDAP configuration.

## Example Usage

```terraform
resource "mongodbatlas_project" "test" {
  name   = "NAME OF THE PROJECT"
  org_id = "ORG ID"
}

resource "mongodbatlas_ldap_configuration" "test" {
  project_id             = mongodbatlas_project.test.id
  authentication_enabled = true
  hostname               = "HOSTNAME"
  port                   = 636
  bind_username          = "USERNAME"
  bind_password          = "PASSWORD"
}
```

## Example Usage of LDAP with user to DN mapping

```terraform
resource "mongodbatlas_project" "test" {
  name   = "NAME OF THE PROJECT"
  org_id = "ORG ID"
}

resource "mongodbatlas_ldap_configuration" "test" {
  project_id             = mongodbatlas_project.test.id
  authentication_enabled = true
  hostname               = "HOSTNAME"
  port                   = 636
  bind_username          = "USERNAME"
  bind_password          = "PASSWORD"
  ca_certificate         = "CA CERTIFICATE"
  authz_query_template   = "{USER}?memberOf?base"
  user_to_dn_mapping {
    match      = "(.+)"
    ldap_query = "DC=example,DC=com??sub?(userPrincipalName={0})"
  }
}
```

## Argument Reference

* `project_id` - (Required) The unique ID for the project to configure LDAP.
* `authentication_enabled` - (Required) Specifies whether user authentication with LDAP is enabled.
* `authorization_enabled` - (Optional) Specifies whether user authorization with LDAP is enabled. You cannot enable user authorization with LDAP without first enabling user authentication with LDAP.
* `hostname` - (Required) The hostname or IP address of the LDAP server. The server must be visible to the internet or connected to your Atlas cluster with VPC Peering.
* `port` - (Optional) The port to which the LDAP server listens for client connections. Default: `636`
* `bind_username` - (Required) The user DN that Atlas uses to connect to the LDAP server. Must be the full DN, such as `CN=BindUser,CN=Users,DC=myldapserver,DC=mycompany,DC=com`.
* `bind_password` - (Required) The password used to authenticate the `bind_username`.
* `ca_certificate` - (Optional) CA certificate used to verify the identify of the LDAP server. Self-signed certificates are allowed.
* `authz_query_template` - (Optional) An LDAP query template that Atlas executes to obtain the LDAP groups to which the authenticated user belongs. Used only for user authorization. Use the {USER} placeholder in the URL to substitute the authenticated username. The query is relative to the host specified with hostname. The formatting for the query must conform to RFC4515 and RFC 4516. If you do not provide a query template, Atlas attempts to use the default value: `{USER}?memberOf?base`.
* `user_to_dn_mapping` - (Optional) Maps an LDAP username for authentication to an LDAP Distinguished Name (DN). Each document contains a `match` regular expression and either a `substitution` or `ldap_query` template used to transform the LDAP username extracted from the regular expression. Atlas steps through the each document in the array in the given order, checking the authentication username against the `match` filter. If a match is found, Atlas applies the transformation and uses the output to authenticate the user. Atlas does not check the remaining documents in the array. For more details and examples see the [MongoDB Atlas API Reference](https://docs.atlas.mongodb.com/reference/api/ldaps-configuration-save/).
* `user_to_dn_mapping.0.match` - (Optional) A regular expression to match against a provided LDAP username. Each parenthesis-enclosed section represents a regular expression capture group used by the `substitution` or `ldap_query` template.
* `user_to_dn_mapping.0.substitution` - (Optional) An LDAP Distinguished Name (DN) formatting template that converts the LDAP name matched by the `match` regular expression into an LDAP Distinguished Name. Each bracket-enclosed numeric value is replaced by the corresponding regular expression capture group extracted from the LDAP username that matched the `match` regular expression.
* `user_to_dn_mapping.0.ldap_query` - (Optional) An LDAP query formatting template that inserts the LDAP name matched by the `match` regular expression into an LDAP query URI as specified by RFC 4515 and RFC 4516. Each numeric value is replaced by the corresponding regular expression capture group extracted from the LDAP username that matched the `match` regular expression.

## Import

LDAP Configuration must be imported using project ID, e.g.

```
$ terraform import mongodbatlas_ldap_configuration.test 5d09d6a59ccf6445652a444a
```

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/ldaps-configuration-save)

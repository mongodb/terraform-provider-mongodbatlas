---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: mongodbatlas_federated_settings_identity_provider"
sidebar_current: "docs-mongodbatlas-federated-settings-identity-provider"
description: |-
    Provides a federated settings Identity Provider resource.
---

# mongodbatlas_federated_settings_identity_provider

`mongodbatlas_federated_settings_identity_provider` provides an Atlas Cloud Federated Settings Identity Provider resource provides a subset of settings to be maintained post import of the existing resource.
## Example Usage

```terraform
resource "mongodbatlas_federated_settings_identity_provider" "identity_provider" {
  federation_settings_id     = "627a9687f7f7f7f774de306f14"
  name = "mongodb_federation_test"
  associated_domains           = ["yourdomain.com"]
  sso_debug_enabled = true
  status = "ACTIVE"
  sso_url = "https://mysso.oktapreview.com/app/mysso_terraformtestsso/exk17q7f7f7f7f50h8/sso/saml"
  issuer_uri = "http://www.okta.com/exk17q7f7f7f7fp50h8"
  request_binding = "HTTP-POST"
  response_signature_algorithm = "SHA-256"
}
```

## Argument Reference

* `federation_settings_id` - (Required) Unique 24-hexadecimal digit string that identifies the federated authentication configuration.
* `name` - (Required) Human-readable label that identifies the identity provider.
* `associated_domains` - (Required) List that contains the domains associated with the identity provider.
* `sso_debug_enabled` - (Required) Flag that indicates whether the identity provider has SSO debug enabled.
* `status`- (Required) String enum that indicates whether the identity provider is active or not. Accepted values are ACTIVE or INACTIVE.
* `issuer_uri` - (Required) Identifier for the issuer of the SAML Assertion.
* `sso_url` - (Required) URL of the receiver of the SAML AuthNRequest.
* `request_binding` - (Required) SAML Authentication Request Protocol binding used to send the AuthNRequest. Atlas supports the following binding values:
    - HTTP POST
    - HTTP REDIRECT
* `response_signature_algorithm` - (Required) Algorithm used to encrypt the IdP signature. Atlas supports the following signature algorithm values:
    - SHA-1


## Attributes Reference

In addition to all arguments above, the following attributes are exported:


### FederatedSettingsIdentityProvider

* `idp_id` - Unique 20-hexadecimal digit string that identifies the IdP.

## Import

Identity Provider **must** be imported before using federation_settings_id-idp_id, e.g.

```
$ terraform import mongodbatlas_federated_settings_identity_provider.identity_provider 6287a663c660f52b1c441c6c-0oad4fas87jL5Xnk1297
```

For more information see: [MongoDB Atlas API Reference.](https://www.mongodb.com/docs/atlas/reference/api/federation-configuration/)
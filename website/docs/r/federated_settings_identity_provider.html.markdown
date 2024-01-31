---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: mongodbatlas_federated_settings_identity_provider"
sidebar_current: "docs-mongodbatlas-federated-settings-identity-provider"
description: |-
    Provides a federated settings Identity Provider resource.
---

# Resource: mongodbatlas_federated_settings_identity_provider

`mongodbatlas_federated_settings_identity_provider` provides an Atlas federated settings identity provider resource provides a subset of settings to be maintained post import of the existing resource.

-> **NOTE:** OIDC Workforce IdP is currently in preview. To learn more about OIDC and existing limitations see the [OIDC Authentication Documentation](https://www.mongodb.com/docs/atlas/security-oidc/).
## Example Usage

~> **IMPORTANT** You **MUST** import this resource before you can manage it with this provider. 

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
* `associated_domains` - List that contains the domains associated with the identity provider.
* `sso_debug_enabled` - Flag that indicates whether the identity provider has SSO debug enabled.
* `status`- String enum that indicates whether the identity provider is active or not. Accepted values are ACTIVE or INACTIVE.
* `issuer_uri` - (Required) Unique string that identifies the issuer of the IdP.
* `sso_url` - Unique string that identifies the intended audience of the SAML assertion.
* `request_binding` - SAML Authentication Request Protocol HTTP method binding (`POST` or `REDIRECT`) that Federated Authentication uses to send the authentication request. Atlas supports the following binding values:
    - HTTP POST
    - HTTP REDIRECT
* `response_signature_algorithm` - Signature algorithm that Federated Authentication uses to encrypt the identity provider signature.  Valid values include `SHA-1 `and `SHA-256`.
* `protocol` - The protocol of the identity provider. Either `SAML` or `OIDC`.
* `audience_claim` - Identifier of the intended recipient of the token.
* `client_id` - Client identifier that is assigned to an application by the Identity Provider.
* `groups_claim` - Identifier of the claim which contains IdP Group IDs in the token.
* `requested_scopes` - Scopes that MongoDB applications will request from the authorization endpoint.
* `user_claim` - Identifier of the claim which contains the user ID in the token.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:


### FederatedSettingsIdentityProvider

* `okta_idp_id` - Unique 20-hexadecimal digit string that identifies the IdP.
* `idp_id` - Unique 24-hexadecimal digit string that identifies the IdP.

## Import

Identity Provider **must** be imported before using federation_settings_id-idp_id, e.g.

```
$ terraform import mongodbatlas_federated_settings_identity_provider.identity_provider 6287a663c660f52b1c441c6c-0oad4fas87jL5Xnk1297
```

**WARNING:** Starting from terraform provider version 1.16.0, to import the resource a 24-hexadecimal digit string that identifies the IdP (`idp_id`) will have to be used instead of `okta_idp_id`. See more [here](../guides/1.15.0-upgrade-guide.html.markdown)

For more information see: [MongoDB Atlas API Reference.](https://www.mongodb.com/docs/atlas/reference/api/federation-configuration/)
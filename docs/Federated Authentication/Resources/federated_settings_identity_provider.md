# Resource: mongodbatlas_federated_settings_identity_provider

`mongodbatlas_federated_settings_identity_provider` provides an Atlas federated settings identity provider resource provides a subset of settings to be maintained post import of the existing resource.

## Example Usage

~> **IMPORTANT** If you want to use a SAML Identity Provider, you **MUST** import this resource before you can manage it with this provider. 

SAML IdP:

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

OIDC IdP:

```
resource "mongodbatlas_federated_settings_identity_provider" "oidc" {
  federation_settings_id = data.mongodbatlas_federated_settings.this.id
  audience               = var.token_audience
  authorization_type     = "USER"
  description            = "oidc"
  issuer_uri = "https://sts.windows.net/${azurerm_user_assigned_identity.this.tenant_id}/"
  idp_type   = "WORKLOAD"
  name       = "OIDC-for-azure"
  protocol   = "OIDC"
  user_claim = "sub"
}
```

## Argument Reference

* `federation_settings_id` - (Required) Unique 24-hexadecimal digit string that identifies the federated authentication configuration.
* `name` - (Required) Human-readable label that identifies the identity provider.
* `description` - (Required for OIDC IdPs) The description of the identity provider.
* `authorization_type` - (Required for OIDC IdPs) Indicates whether authorization is granted based on group membership or user ID. Valid values are `GROUP` or `USER`.
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
* `audience` - (Required for OIDC IdPs) Identifier of the intended recipient of the token used in OIDC IdP.
* `client_id` - Client identifier that is assigned to an application by the OIDC Identity Provider.
* `groups_claim` - (Required for OIDC IdP with `authorization_type = GROUP`) Identifier of the claim which contains OIDC IdP Group IDs in the token.
* `requested_scopes` - Scopes that MongoDB applications will request from the authorization endpoint used for OIDC IdPs.
* `user_claim` - (Required for OIDC IdP) Identifier of the claim which contains the user ID in the token used for OIDC IdPs.
userClaim is required for OIDC IdP with authorizationType GROUP and USER.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:


### FederatedSettingsIdentityProvider

* `okta_idp_id` - Unique 20-hexadecimal digit string that identifies the IdP.
* `idp_id` - Unique 24-hexadecimal digit string that identifies the IdP.

## Import

Identity Provider **must** be imported before using federation_settings_id-idp_id, e.g.

```
$ terraform import mongodbatlas_federated_settings_identity_provider.identity_provider 6287a663c660f52b1c441c6c-0oad4fas87jL5Xnk12971234
```

For more information see: [MongoDB Atlas API Reference.](https://www.mongodb.com/docs/atlas/reference/api/federation-configuration/)
# Data Source: mongodbatlas_federated_settings_identity_provider

`mongodbatlas_federated_settings_identity_provider` provides a federated settings identity provider data source. Atlas federated settings identity provider provides federated settings outputs for the configured identity provider.

## Example Usage

```terraform
resource "mongodbatlas_federated_settings_identity_provider" "identity_provider" {
  federation_settings_id     = "627a9687f7f7f7f774de306f14"
  name = "mongodb_federation_test"
  associated_domains           = ["yourdomain.com"]
  sso_debug_enabled = true
  status = "ACTIVE"
  sso_url = "https://mysso.oktapreview.com/app/mysso_terraformtest_1/exk177f7f7f70h8/sso/saml"
  issuer_uri = "http://www.okta.com/exk17f7f7f7f7p50h8"
  request_binding = "HTTP-POST"
  response_signature_algorithm = "SHA-256"
}

data "mongodbatlas_federated_settings_identity_provider" "identity_provider_ds" {
  federation_settings_id = mongodbatlas_federated_settings_identity_provider.identity_provider.id
  identity_provider_id   = "0oad4fas87jL5Xnk12971234"
}

```

## Argument Reference

* `federation_settings_id` - (Required) Unique 24-hexadecimal digit string that identifies the federated authentication configuration.
* `identity_provider_id` - (Required) Unique 24-hexadecimal digit string that identifies the IdP. 

## Attributes Reference

In addition to all arguments above, the following attributes are exported:


### FederatedSettingsIdentityProvider

* `description` - The description of the identity provider.
* `authorization_type` - Indicates whether authorization is granted based on group membership or user ID. Valid values are `GROUP` or `USER`.
* `acs_url` - Assertion consumer service URL to which the IdP sends the SAML response.
* `associated_domains` - List that contains the configured domains from which users can log in for this IdP.
* `associated_orgs` - List that contains the organizations from which users can log in for this IdP.
* `domain_allow_list` - List that contains the approved domains from which organization users can log in.
* `domain_restriction_enabled` - Flag that indicates whether domain restriction is enabled for the connected organization.
* `org_id` - Unique 24-hexadecimal digit string that identifies the organization that contains your projects.
* `post_auth_role_grants` - List that contains the default roles granted to users who authenticate through the IdP in a connected organization. If you provide a postAuthRoleGrants field in the request, the array that you provide replaces the current postAuthRoleGrants.
* `idp_id` - Unique 24-hexadecimal digit string that identifies the IdP
* `protocol` - The protocol of the identity provider. Valid values are `SAML` or `OIDC`.
* `audience` - Identifier of the intended recipient of the token.
* `client_id` - Client identifier that is assigned to an application by the Identity Provider.
* `groups_claim` - Identifier of the claim which contains IdP Group IDs in the token.
* `requested_scopes` - Scopes that MongoDB applications will request from the authorization endpoint.
* `user_claim` - Identifier of the claim which contains the user ID in the token.
* `idp_type` - Type of the identity provider. Valid values are `WORKFORCE` or `WORKLOAD`.

  ### Role_mappings
* `external_group_name` - Unique human-readable label that identifies the identity provider group to which this role mapping applies.
* `id` - Unique 24-hexadecimal digit string that identifies this role mapping.
* `role_assignments` - Atlas roles and the unique identifiers of the groups and organizations associated with each role.
* `group_id` - Unique identifier of the project to which you want the role mapping to apply.
* `role` - Specifies the Role that is attached to the Role Mapping.
### User Conflicts
* `email_address` - Email address of the the user that conflicts with selected domains.
* `federation_settings_id` - Unique 24-hexadecimal digit string that identifies the federated authentication configuration.
* `first_name` - First name of the the user that conflicts with selected domains.
* `last_name` - Last name of the the user that conflicts with selected domains.
* `user_id` - Name of the Atlas user that conflicts with selected domains.
* `audience_uri` - Identifier for the intended audience of the SAML Assertion.
* `display_name` - Human-readable label that identifies the IdP.
* `issuer_uri` - Identifier for the issuer of the SAML Assertion.
* `okta_idp_id` - Unique 20-hexadecimal digit string that identifies the IdP.
### Pem File Info - List that contains the file information, including: start date, and expiration date for the identity provider's PEM-encoded public key certificate.
* `not_after` - Expiration  Date.
* `not_before` - Start Date.
* `file_name` - Filename of certificate
* `request_binding` - SAML Authentication Request Protocol binding used to send the AuthNRequest. Atlas supports the following binding values:
    - HTTP POST
    - HTTP REDIRECT
* `response_signature_algorithm` - Algorithm used to encrypt the IdP signature. Atlas supports the following signature algorithm values:
    - SHA-1
    - SHA-256
* `sso_debug_enabled` - Flag that indicates whether the IdP has enabled Bypass SAML Mode. Enabling this mode generates a URL that allows you bypass SAML and login to your organizations at any point. You can authenticate with this special URL only when Bypass Mode is enabled. Set this parameter to true during testing. This keeps you from getting locked out of MongoDB.
* `sso_url` - URL of the receiver of the SAML AuthNRequest.
* `status` - Label that indicates whether the identity provider is active. The IdP is Inactive until you map at least one domain to the IdP.


For more information see: [MongoDB Atlas API Reference.](https://www.mongodb.com/docs/atlas/reference/api/federation-configuration/)

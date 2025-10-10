---
page_title: "Migration Guide: Service Accounts"
---

# Migration Guide: Service Accounts

This guide helps you migrate from Programmatic Access Key (PAK) authentication to Service
Accounts (SA) authentication without impacting your deployment. SA simplifies authentication by eliminating the need to create new Atlas-specific user identities and permission credentials.

For more information on SA, see [Service Accounts Overview](https://www.mongodb.com/docs/atlas/api/service-accounts-overview/) in the MongoDB documentation.
These are the [MongoDB Atlas Service Account Limits](https://www.mongodb.com/docs/manual/reference/limits/#mongodb-atlas-service-account-limits).

For more comprehensive information about the authentication methods available in the MongoDB Atlas Provider, see the [Authenticate the Provider](../#authenticate-the-provider) section in the provider's main page.

**NOTE:** Service Accounts can't be used with `mongodbatlas_event_trigger` resources as its API doesn't support it yet.

## Before You Begin

See [Get Started with the Atlas Administration API](https://www.mongodb.com/docs/atlas/configure-api-access/#grant-programmatic-access-to-an-organization) in the MongoDB Atlas documentation for detailed instructions on how to create your organization's SA and granting the required access.

## PAK in environment variables

If your PAK is configured using environment variables, for example:
```shell
$ export MONGODB_ATLAS_PUBLIC_API_KEY="<ATLAS_PUBLIC_API_KEY>"
$ export MONGODB_ATLAS_PRIVATE_API_KEY="<ATLAS_PRIVATE_API_KEY>"
```

You need to stop setting those environment variables and set the environment variables for the SA you created, e.g.:
```shell
$ export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
$ export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

Nothing else is required, you can run `terraform plan` to check that everything is working fine.

## PAK in provider attributes

If your PAK is in provider attributes, for example:
```terraform
provider "mongodbatlas" {
  public_key = var.mongodbatlas_public_key
  private_key = var.mongodbatlas_private_key
}
```

You need to remove those attributes from the provider configuration and add the SA attributes along with defining new variables for them, e.g.:

```terraform
provider "mongodbatlas" {
  client_id = var.mongodbatlas_client_id
  client_secret = var.mongodbatlas_client_secret
}
```

Nothing else is required, you can run `terraform plan` to check that everything is working fine.

## PAK in AWS Secrets Manager

If the PAK credentials are stored in a secret in AWS Secrets Manager, you have to remove the raw values for `public_key` and `private_key`, and add the raw values for `client_id` and `client_secret`.

## Service Account Token

Instead of using your Client ID and Client Secret, you can alternatively generate the SA Token yourself, see the [Authenticate the Provider](../#authenticate-the-provider) section in the provider's main page for more information.

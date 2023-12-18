---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: x509_authentication_database_user"
sidebar_current: "docs-mongodbatlas-resource-x509-authentication-database-user"
description: |-
    Provides a X509 Authentication Database User resource.
---

# Resource: mongodbatlas_x509_authentication_database_user

`mongodbatlas_x509_authentication_database_user` provides a X509 Authentication Database User resource. The mongodbatlas_x509_authentication_database_user resource lets you manage MongoDB users who authenticate using X.509 certificates. You can manage these X.509 certificates or let Atlas do it for you.

| Management  | Description  |
|---|---|
| Atlas  | Atlas manages your Certificate Authority and can generate certificates for your MongoDB users. No additional X.509 configuration is required.  |
| Customer  |  You must provide a Certificate Authority and generate certificates for your MongoDB users. |

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

-> **NOTE:** Before provider version 1.14.0, Self-managed X.509 Authentication was disabled for the project when this resource was deleted. Starting from that version onward, it will not be disabled, allowing other users to continue using X.509 within the same project.

## Example Usages

### Example Usage: Generate an Atlas-managed X.509 certificate for a MongoDB user
```terraform
resource "mongodbatlas_database_user" "user" {
  project_id    = "64b926dd56206839b1c8bae9"
  username      = "myUsername"
  x509_type     = "MANAGED"
  database_name = "$external"

  roles {
    role_name     = "atlasAdmin"
    database_name = "admin"
  }

  labels {
    key   = "My Key"
    value = "My Value"
  }
}

resource "mongodbatlas_x509_authentication_database_user" "test" {
  project_id              = mongodbatlas_database_user.user.project_id
  username                = mongodbatlas_database_user.user.username
  months_until_expiration = 2
}
```

### Example Usage: Save a self-managed X.509 certificate for an Atlas project and use it with a dababase user
```terraform
resource "mongodbatlas_x509_authentication_database_user" "test" {
  project_id        = "<PROJECT-ID>"
  customer_x509_cas = <<-EOT
    -----BEGIN CERTIFICATE-----
    MIICmTCCAgICCQDZnHzklxsT9TANBgkqhkiG9w0BAQsFADCBkDELMAkGA1UEBhMC
    VVMxDjAMBgNVBAgMBVRleGFzMQ8wDQYDVQQHDAZBdXN0aW4xETAPBgNVBAoMCHRl
    c3QuY29tMQ0wCwYDVQQLDARUZXN0MREwDwYDVQQDDAh0ZXN0LmNvbTErMCkGCSqG
    SIb3DQEJARYcbWVsaXNzYS5wbHVua2V0dEBtb25nb2RiLmNvbTAeFw0yMDAyMDQy
    MDQ2MDFaFw0yMTAyMDMyMDQ2MDFaMIGQMQswCQYDVQQGEwJVUzEOMAwGA1UECAwF
    VGV4YXMxDzANBgNVBAcMBkF1c3RpbjERMA8GA1UECgwIdGVzdC5jb20xDTALBgNV
    BAsMBFRlc3QxETAPBgNVBAMMCHRlc3QuY29tMSswKQYJKoZIhvcNAQkBFhxtZWxp
    c3NhLnBsdW5rZXR0QG1vbmdvZGIuY29tMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCB
    iQKBgQCf1LRqr1zftzdYx2Aj9G76tb0noMPtj6faGLlPji1+m6Rn7RWD9L0ntWAr
    cURxvypa9jZ9MXFzDtLevvd3tHEmfrUT3ukNDX6+Jtc4kWm+Dh2A70Pd+deKZ2/O
    Fh8audEKAESGXnTbeJCeQa1XKlIkjqQHBNwES5h1b9vJtFoLJwIDAQABMA0GCSqG
    SIb3DQEBCwUAA4GBADMUncjEPV/MiZUcVNGmktP6BPmEqMXQWUDpdGW2+Tg2JtUA
    7MMILtepBkFzLO+GlpZxeAlXO0wxiNgEmCRONgh4+t2w3e7a8GFijYQ99FHrAC5A
    iul59bdl18gVqXia1Yeq/iK7Ohfy/Jwd7Hsm530elwkM/ZEkYDjBlZSXYdyz
    -----END CERTIFICATE-----"
  EOT
}

resource "mongodbatlas_database_user" "user" {
  project_id    = "64b926dd56206839b1c8bae9"
  username      = "myUsername"
  x509_type     = "CUSTOMER" # Make sure to set x509_type = "CUSTOMER"
  database_name = "$external"

  roles {
    role_name     = "atlasAdmin"
    database_name = "admin"
  }

  labels {
    key   = "My Key"
    value = "My Value"
  }
}
```

## Argument Reference

* `project_id` - (Required) Identifier for the Atlas project associated with the X.509 configuration.
* `months_until_expiration` - (Required) A number of months that the created certificate is valid for before expiry, up to 24 months. By default is 3.
* `username` - (Optional) Username of the database user to create a certificate for.
* `customer_x509_cas` - (Optional) PEM string containing one or more customer CAs for database user authentication.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:

* `current_certificate` - Contains the last X.509 certificate and private key created for a database user.

  #### Certificates
* `certificates` - Array of objects where each details one unexpired database user certificate.

* `certificates.#.id` - Serial number of this certificate.
* `certificates.#.created_at` - Timestamp in ISO 8601 date and time format in UTC when Atlas created this X.509 certificate.
* `certificates.#.group_id` - Unique identifier of the Atlas project to which this certificate belongs.
* `certificates.#.not_after` - Timestamp in ISO 8601 date and time format in UTC when this certificate expires.
* `certificates.#.subject` - Fully distinguished name of the database user to which this certificate belongs. To learn more, see RFC 2253.

## Import

X.509 Certificates for a User can be imported using project ID and username, in the format `project_id-username`, e.g.

```
$ terraform import mongodbatlas_x509_authentication_database_user.test 1112222b3bf99403840e8934-myUsername
```

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/x509-configuration-get-certificates/)


Current X.509 Configuration can be imported using project ID, in the format `project_id`, e.g.

```
$ terraform import mongodbatlas_x509_authentication_database_user.test 1112222b3bf99403840e8934
```

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/x509-configuration-get-certificates/)
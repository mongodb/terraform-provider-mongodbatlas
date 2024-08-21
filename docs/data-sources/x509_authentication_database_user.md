# Data Source: mongodbatlas_x509_authentication_database_user

`mongodbatlas_x509_authentication_database_user` describes a X509 Authentication Database User. This represents a X509 Authentication Database User.

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

## Example Usages

### Example Usage: Generate an Atlas-managed X.509 certificate for a MongoDB user
```terraform
resource "mongodbatlas_database_user" "user" {
  project_id    = "<PROJECT-ID>"
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

data "mongodbatlas_x509_authentication_database_user" "test" {
  project_id = mongodbatlas_x509_authentication_database_user.test.project_id
  username   = mongodbatlas_x509_authentication_database_user.test.username
}
```

### Example Usage: Save a customer-managed X.509 configuration for an Atlas project
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

data "mongodbatlas_x509_authentication_database_user" "test" {
  project_id = mongodbatlas_x509_authentication_database_user.test.project_id
}
```

## Argument Reference

* `project_id` - (Required) Identifier for the Atlas project associated with the X.509 configuration.
* `username` - (Optional) Username of the database user to create a certificate for.

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


See [MongoDB Atlas - X509 User Certificates](https://docs.atlas.mongodb.com/reference/api/x509-configuration-get-certificates/) and [MongoDB Atlas - Current X509 Configuratuion](https://docs.atlas.mongodb.com/reference/api/x509-configuration-get-current/) Documentation for more information.
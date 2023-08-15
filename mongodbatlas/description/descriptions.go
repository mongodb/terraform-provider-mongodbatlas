// Copyright 2023 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//nolint

package description

const (
	DatabaseUserResource = "`mongodbatlas_database_user` provides a Database User resource. This represents a database user which will be applied to all clusters within the project.\n\nEach user has a set of roles that provide access to the project’s databases. User's roles apply to all the clusters in the project: if two clusters have a `products` database and a user has a role granting `read` access on the products database, the user has that access on both clusters."
	DatabaseUserDS       = "`mongodbatlas_database_user` describe a Database User. This represents a database user which will be applied to all clusters within the project.\n\nEach user has a set of roles that provide access to the project’s databases. User's roles apply to all the clusters in the project: if two clusters have a `products` database and a user has a role granting `read` access on the products database, the user has that access on both clusters."
	DatabaseUsersDS      = "`mongodbatlas_database_users` describe all Database Users. This represents a database user which will be applied to all clusters within the project.\n\nEach user has a set of roles that provide access to the project’s databases. User's roles apply to all the clusters in the project: if two clusters have a `products` database and a user has a role granting `read` access on the products database, the user has that access on both clusters."
	ID                   = "ID used during the resource import."
	ProjectID            = "The unique ID for the project to create the database user."
	Username             = "Username for authenticating to MongoDB. USER_ARN or ROLE_ARN if `aws_iam_type` is USER or ROLE."
	Password             = `User's initial password. A value is required to create the database user, however the argument but may be removed from your Terraform configuration after user creation without impacting the user, password or Terraform management. //nolint:gosec // This is just a message not a credential
	IMPORTANT --- Passwords may show up in Terraform related logs and it will be stored in the Terraform state file as plain-text. Password can be changed after creation using your preferred method, e.g. via the MongoDB Atlas UI, to ensure security.  If you do change management of the password to outside of Terraform be sure to remove the argument from the Terraform configuration so it is not inadvertently updated to the original password.`
	X509Type     = "X.509 method by which the provided username is authenticated. If no value is given, Atlas uses the default value of NONE. The accepted types are:\n* `NONE` -	The user does not use X.509 authentication.\n* `MANAGED` - The user is being created for use with Atlas-managed X.509.Externally authenticated users can only be created on the `$external` database.\n* `CUSTOMER` -  The user is being created for use with Self-Managed X.509. Users created with this x509Type require a Common Name (CN) in the username field. Externally authenticated users can only be created on the `$external` database."
	AWSIAMType   = "If this value is set, the new database user authenticates with AWS IAM credentials. If no value is given, Atlas uses the default value of NONE. The accepted types are:\n* `NONE` -	The user does not use AWS IAM credentials.* `USER` - New database user has AWS IAM user credentials.\n* `ROLE` -  New database user has credentials associated with an AWS IAM role."
	LDAPAuthYype = "Method by which the provided `username` is authenticated. If no value is given, Atlas uses the default value of `NONE`.\n* `NONE` -	Atlas authenticates this user through [SCRAM-SHA](https://docs.mongodb.com/manual/core/security-scram/), not LDAP.\n* `USER` - LDAP server authenticates this user through the user's LDAP user. `username` must also be a fully qualified distinguished name, as defined in [RFC-2253](https://tools.ietf.org/html/rfc2253).\n* `GROUP` - LDAP server authenticates this user using their LDAP user and authorizes this user using their LDAP group. To learn more about LDAP security, see [Set up User Authentication and Authorization with LDAP](https://docs.atlas.mongodb.com/security-ldaps). `username` must also be a fully qualified distinguished name, as defined in [RFC-2253](https://tools.ietf.org/html/rfc2253)."
	Roles        = `Block mapping a user's role to a database / collection. A role allows the user to perform particular actions on the specified database. A role on the admin database can include privileges that apply to the other databases as well.

	-> **NOTE:** The available privilege actions for custom MongoDB roles support a subset of MongoDB commands. See Unsupported Commands in M10+ Clusters for more information.
	
	~> **IMPORTANT:** If a user is assigned a custom MongoDB role, they cannot be assigned any other roles.`
	RoleName       = "Name of the role to grant. See [Create a Database User](https://docs.atlas.mongodb.com/reference/api/database-users-create-a-user/) `roles.roleName` for valid values and restrictions."
	DatabaseName   = "Database on which the user has the specified role. A role on the `admin` database can include privileges that apply to the other databases."
	CollectionName = "Collection for which the role applies. You can specify a collection for the `read` and `readWrite` roles. If you do not specify a collection for `read` and `readWrite`, the role applies to all collections in the database (excluding some collections in the `system`. database)."
	Labels         = "Containing key-value pairs that tag and categorize the database user. Each key and value has a maximum length of 255 characters."
	Key            = "The key that you want to write."
	Value          = "The value that you want to write."
	Scopes         = "Array of clusters and Atlas Data Lakes that this user has access to. If omitted, Atlas grants the user access to all the clusters and Atlas Data Lakes in the project by default."
	Name           = "Name of the cluster or Atlas Data Lake that the user has access to."
	Type           = "Type of resource that the user has access to. Valid values are: `CLUSTER` and `DATA_LAKE`"
)

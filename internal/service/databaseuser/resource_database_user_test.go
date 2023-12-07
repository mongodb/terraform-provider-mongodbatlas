package databaseuser_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"go.mongodb.org/atlas-sdk/v20231115002/admin"
)

func TestAccConfigRSDatabaseUser_basic(t *testing.T) {
	var (
		dbUser       admin.CloudDatabaseUser
		resourceName = "mongodbatlas_database_user.basic_ds"
		username     = acctest.RandomWithPrefix("dbUser")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyDatabaseUser,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigDatabaseUserBasic(projectName, orgID, "atlasAdmin", username, "First Key", "First value"),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckDatabaseUserExists(resourceName, &dbUser),
					acc.CheckDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", "test-acc-password"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "1"),
				),
			},
			{
				Config: acc.ConfigDatabaseUserBasic(projectName, orgID, "read", username, "Second Key", "Second value"),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckDatabaseUserExists(resourceName, &dbUser),
					acc.CheckDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", "test-acc-password"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "1"),
				),
			},
		},
	})
}

func TestAccConfigRSDatabaseUser_withX509TypeCustomer(t *testing.T) {
	var (
		dbUser       admin.CloudDatabaseUser
		resourceName = "mongodbatlas_database_user.test"
		username     = "CN=ellen@example.com,OU=users,DC=example,DC=com"
		x509Type     = "CUSTOMER"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyDatabaseUser,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigDatabaseUserWithX509Type(projectName, orgID, "atlasAdmin", username, "First Key", "First value", x509Type),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckDatabaseUserExists(resourceName, &dbUser),
					acc.CheckDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "x509_type", x509Type),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "$external"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
				),
			},
		},
	})
}

func TestAccConfigRSDatabaseUser_withX509TypeManaged(t *testing.T) {
	var (
		dbUser       admin.CloudDatabaseUser
		resourceName = "mongodbatlas_database_user.test"
		username     = acctest.RandomWithPrefix("test-acc")
		x509Type     = "MANAGED"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyDatabaseUser,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigDatabaseUserWithX509Type(projectName, orgID, "atlasAdmin", username, "First Key", "First value", x509Type),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckDatabaseUserExists(resourceName, &dbUser),
					acc.CheckDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "x509_type", x509Type),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "$external"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
				),
			},
		},
	})
}

func TestAccConfigRSDatabaseUser_withAWSIAMType(t *testing.T) {
	var (
		dbUser       admin.CloudDatabaseUser
		resourceName = "mongodbatlas_database_user.test"
		username     = "arn:aws:iam::358363220050:user/mongodb-aws-iam-auth-test-user"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyDatabaseUser,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigDatabaseUserWithAWSIAMType(projectName, orgID, "atlasAdmin", username, "First Key", "First value"),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckDatabaseUserExists(resourceName, &dbUser),
					acc.CheckDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "aws_iam_type", "USER"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "$external"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
				),
			},
		},
	})
}

func TestAccConfigRSDatabaseUser_withAWSIAMType_import(t *testing.T) {
	var (
		dbUser       admin.CloudDatabaseUser
		resourceName = "mongodbatlas_database_user.test"
		username     = os.Getenv("TEST_DB_USER_IAM_ARN")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	if username == "" {
		username = "arn:aws:iam::358363220050:user/mongodb-aws-iam-auth-test-user"
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyDatabaseUser,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigDatabaseUserWithAWSIAMType(projectName, orgID, "atlasAdmin", username, "First Key", "First value"),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckDatabaseUserExists(resourceName, &dbUser),
					acc.CheckDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "aws_iam_type", "USER"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "$external"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasDatabaseUserImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccConfigRSDatabaseUser_WithLabels(t *testing.T) {
	var (
		dbUser       admin.CloudDatabaseUser
		resourceName = "mongodbatlas_database_user.test"
		username     = acctest.RandomWithPrefix("test-acc")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyDatabaseUser,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigDatabaseUserWithLabels(projectName, orgID, "atlasAdmin", username, []admin.ComponentLabel{}),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckDatabaseUserExists(resourceName, &dbUser),
					acc.CheckDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", "test-acc-password"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "0"),
				),
			},
			{
				Config: acc.ConfigDatabaseUserWithLabels(projectName, orgID, "atlasAdmin", username,
					[]admin.ComponentLabel{
						{
							Key:   conversion.StringPtr("key 1"),
							Value: conversion.StringPtr("value 1"),
						},
						{
							Key:   conversion.StringPtr("key 2"),
							Value: conversion.StringPtr("value 2"),
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckDatabaseUserExists(resourceName, &dbUser),
					acc.CheckDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", "test-acc-password"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "2"),
				),
			},
			{
				Config: acc.ConfigDatabaseUserWithLabels(projectName, orgID, "read", username,
					[]admin.ComponentLabel{
						{
							Key:   conversion.StringPtr("key 4"),
							Value: conversion.StringPtr("value 4"),
						},
						{
							Key:   conversion.StringPtr("key 3"),
							Value: conversion.StringPtr("value 3"),
						},
						{
							Key:   conversion.StringPtr("key 2"),
							Value: conversion.StringPtr("value 2"),
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckDatabaseUserExists(resourceName, &dbUser),
					acc.CheckDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", "test-acc-password"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "3"),
				),
			},
		},
	})
}

func TestAccConfigRSDatabaseUser_withRoles(t *testing.T) {
	var (
		dbUser       admin.CloudDatabaseUser
		resourceName = "mongodbatlas_database_user.test"
		username     = acctest.RandomWithPrefix("test-acc-user-")
		password     = acctest.RandomWithPrefix("test-acc-pass-")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyDatabaseUser,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigDatabaseUserWithRoles(username, password, projectName, orgID,
					[]*admin.DatabaseUserRole{
						{
							RoleName:       "read",
							DatabaseName:   "admin",
							CollectionName: conversion.StringPtr("stir"),
						},
						{
							RoleName:       "read",
							DatabaseName:   "admin",
							CollectionName: conversion.StringPtr("unpledged"),
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckDatabaseUserExists(resourceName, &dbUser),
					acc.CheckDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "roles.0.collection_name", "stir"),
					resource.TestCheckResourceAttr(resourceName, "roles.1.collection_name", "unpledged"),
				),
			},
			{
				Config: acc.ConfigDatabaseUserWithRoles(username, password, projectName, orgID,
					[]*admin.DatabaseUserRole{
						{
							RoleName:     "read",
							DatabaseName: "admin",
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckDatabaseUserExists(resourceName, &dbUser),
					acc.CheckDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "1"),
				),
			},
		},
	})
}

func TestAccConfigRSDatabaseUser_withScopes(t *testing.T) {
	var (
		dbUser       admin.CloudDatabaseUser
		resourceName = "mongodbatlas_database_user.test"
		username     = acctest.RandomWithPrefix("test-acc-user-")
		password     = acctest.RandomWithPrefix("test-acc-pass-")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		clusterName  = acctest.RandomWithPrefix("test-acc-cluster")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyDatabaseUser,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigDatabaseUserWithScopes(username, password, projectName, orgID, "atlasAdmin", clusterName,
					[]*admin.UserScope{
						{
							Name: "test-acc-nurk4llu2z",
							Type: "CLUSTER",
						},
						{
							Name: "test-acc-nurk4llu2z",
							Type: "DATA_LAKE",
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckDatabaseUserExists(resourceName, &dbUser),
					acc.CheckDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "scopes.#", "2"),
				),
			},
			{
				Config: acc.ConfigDatabaseUserWithScopes(username, password, projectName, orgID, "atlasAdmin", clusterName,
					[]*admin.UserScope{
						{
							Name: "test-acc-nurk4llu2z",
							Type: "CLUSTER",
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckDatabaseUserExists(resourceName, &dbUser),
					acc.CheckDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "scopes.#", "1"),
				),
			},
		},
	})
}

func TestAccConfigRSDatabaseUser_withScopesAndEmpty(t *testing.T) {
	var (
		dbUser       admin.CloudDatabaseUser
		resourceName = "mongodbatlas_database_user.test"
		username     = acctest.RandomWithPrefix("test-acc-user-")
		password     = acctest.RandomWithPrefix("test-acc-pass-")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		clusterName  = acctest.RandomWithPrefix("test-acc-cluster")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyDatabaseUser,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigDatabaseUserWithScopes(username, password, projectName, orgID, "atlasAdmin", clusterName,
					[]*admin.UserScope{
						{
							Name: "test-acc-nurk4llu2z",
							Type: "CLUSTER",
						},
						{
							Name: "test-acc-nurk4llu2z",
							Type: "DATA_LAKE",
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckDatabaseUserExists(resourceName, &dbUser),
					acc.CheckDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "scopes.#", "2"),
				),
			},
			{
				Config: acc.ConfigDatabaseUserWithScopes(username, password, projectName, orgID, "atlasAdmin", clusterName,
					[]*admin.UserScope{},
				),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckDatabaseUserExists(resourceName, &dbUser),
					acc.CheckDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "scopes.#", "0"),
				),
			},
		},
	})
}

func TestAccConfigRSDatabaseUser_withLDAPAuthType(t *testing.T) {
	var (
		dbUser       admin.CloudDatabaseUser
		resourceName = "mongodbatlas_database_user.test"
		username     = "CN=david@example.com,OU=users,DC=example,DC=com"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyDatabaseUser,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigDatabaseUserWithLDAPAuthType(projectName, orgID, "atlasAdmin", username, "First Key", "First value"),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckDatabaseUserExists(resourceName, &dbUser),
					acc.CheckDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "ldap_auth_type", "USER"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "$external"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
				),
			},
		},
	})
}

func TestAccConfigRSDatabaseUser_importBasic(t *testing.T) {
	var (
		username     = fmt.Sprintf("test-username-%s", acctest.RandString(5))
		resourceName = "mongodbatlas_database_user.basic_ds"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyDatabaseUser,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigDatabaseUserBasic(projectName, orgID, "read", username, "First Key", "First value"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", "test-acc-password"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasDatabaseUserImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccConfigRSDatabaseUser_importX509TypeCustomer(t *testing.T) {
	var (
		dbUser       admin.CloudDatabaseUser
		resourceName = "mongodbatlas_database_user.test"
		username     = "CN=ellen@example.com,OU=users,DC=example,DC=com"
		x509Type     = "CUSTOMER"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyDatabaseUser,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigDatabaseUserWithX509Type(projectName, orgID, "atlasAdmin", username, "First Key", "First value", x509Type),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckDatabaseUserExists(resourceName, &dbUser),
					acc.CheckDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "x509_type", x509Type),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "$external"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasDatabaseUserImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccConfigRSDatabaseUser_importLDAPAuthType(t *testing.T) {
	var (
		dbUser       admin.CloudDatabaseUser
		resourceName = "mongodbatlas_database_user.test"
		username     = "CN=david@example.com,OU=users,DC=example,DC=com"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyDatabaseUser,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigDatabaseUserWithLDAPAuthType(projectName, orgID, "atlasAdmin", username, "First Key", "First value"),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckDatabaseUserExists(resourceName, &dbUser),
					acc.CheckDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "ldap_auth_type", "USER"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "$external"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasDatabaseUserImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccCheckMongoDBAtlasDatabaseUserImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s-%s-%s", ids["project_id"], ids["username"], ids["auth_database_name"]), nil
	}
}

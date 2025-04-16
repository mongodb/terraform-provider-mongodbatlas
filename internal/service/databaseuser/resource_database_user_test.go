package databaseuser_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/databaseuser"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"go.mongodb.org/atlas-sdk/v20250312002/admin"
)

const (
	resourceName         = "mongodbatlas_database_user.test"
	dataSourceName       = "data.mongodbatlas_database_user.test"
	dataSourcePluralName = "data.mongodbatlas_database_users.test"
)

func TestAccConfigRSDatabaseUser_basic(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		username  = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigDatabaseUserBasic(projectID, username, "atlasAdmin", "First Key", "First value"),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", "test-acc-password"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "labels.0.key", "First Key"),
					resource.TestCheckResourceAttr(resourceName, "labels.0.value", "First value"),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "roles.0.role_name", "atlasAdmin"),
					resource.TestCheckNoResourceAttr(resourceName, "description"),
				),
			},
			{
				Config: acc.ConfigDatabaseUserBasic(projectID, username, "read", "Second Key", "Second value"),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", "test-acc-password"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "labels.0.key", "Second Key"),
					resource.TestCheckResourceAttr(resourceName, "labels.0.value", "Second value"),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "roles.0.role_name", "read"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccConfigRSDatabaseUser_withX509TypeCustomer(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		username  = acc.RandomLDAPName()
		x509Type  = "CUSTOMER"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigDatabaseUserWithX509Type(projectID, username, x509Type, "atlasAdmin", "First Key", "First value"),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "x509_type", x509Type),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "$external"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccConfigRSDatabaseUser_withX509TypeManaged(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		x509Type  = "MANAGED"
		username  = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigDatabaseUserWithX509Type(projectID, username, x509Type, "atlasAdmin", "First Key", "First value"),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
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
		projectID = acc.ProjectIDExecution(t)
		username  = acc.RandomIAMUser()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigDatabaseUserWithAWSIAMType(projectID, username, "atlasAdmin", "First Key", "First value"),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "aws_iam_type", "USER"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "$external"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccConfigRSDatabaseUser_withLabelsAndDescription(t *testing.T) {
	var (
		projectID    = acc.ProjectIDExecution(t)
		username     = acc.RandomName()
		description1 = "desc 1"
		description2 = "desc 2"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigDatabaseUserWithLabels(projectID, username, "atlasAdmin", description1, nil),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", "test-acc-password"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "description", description1),
				),
			},
			{
				Config: acc.ConfigDatabaseUserWithLabels(projectID, username, "atlasAdmin", description2,
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
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "description", description2),
					resource.TestCheckResourceAttr(resourceName, "password", "test-acc-password"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "2"),
				),
			},
			{
				Config: acc.ConfigDatabaseUserWithLabels(projectID, username, "read", "",
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
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", "test-acc-password"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "3"),
					resource.TestCheckNoResourceAttr(resourceName, "description"),
				),
			},
		},
	})
}

func TestAccConfigRSDatabaseUser_withRoles(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		username  = acc.RandomName()
		password  = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigDatabaseUserWithRoles(projectID, username, password,
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
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "roles.0.collection_name", "stir"),
					resource.TestCheckResourceAttr(resourceName, "roles.1.collection_name", "unpledged"),
				),
			},
			{
				Config: acc.ConfigDatabaseUserWithRoles(projectID, username, password,
					[]*admin.DatabaseUserRole{
						{
							RoleName:     "read",
							DatabaseName: "admin",
						},
					},
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
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
		projectID     = acc.ProjectIDExecution(t)
		userScopeName = acc.RandomName()
		username      = acc.RandomName()
		password      = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigDatabaseUserWithScopes(projectID, username, password, "atlasAdmin",
					[]*admin.UserScope{
						{
							Name: userScopeName,
							Type: "CLUSTER",
						},
						{
							Name: userScopeName,
							Type: "DATA_LAKE",
						},
					},
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "scopes.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "scopes.0.name", userScopeName),
					resource.TestCheckResourceAttr(resourceName, "scopes.0.type", "CLUSTER"),
					resource.TestCheckResourceAttr(resourceName, "scopes.1.name", userScopeName),
					resource.TestCheckResourceAttr(resourceName, "scopes.1.type", "DATA_LAKE"),
				),
			},
			{
				Config: acc.ConfigDatabaseUserWithScopes(projectID, username, password, "atlasAdmin",
					[]*admin.UserScope{
						{
							Name: userScopeName,
							Type: "CLUSTER",
						},
					},
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "scopes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scopes.0.name", userScopeName),
					resource.TestCheckResourceAttr(resourceName, "scopes.0.type", "CLUSTER"),
				),
			},
		},
	})
}

func TestAccConfigRSDatabaseUser_updateToEmptyScopes(t *testing.T) {
	var (
		projectID     = acc.ProjectIDExecution(t)
		userScopeName = acc.RandomName()
		username      = acc.RandomName()
		password      = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigDatabaseUserWithScopes(projectID, username, password, "atlasAdmin",
					[]*admin.UserScope{
						{
							Name: userScopeName,
							Type: "CLUSTER",
						},
						{
							Name: userScopeName,
							Type: "DATA_LAKE",
						},
					},
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "scopes.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "scopes.0.name", userScopeName),
					resource.TestCheckResourceAttr(resourceName, "scopes.0.type", "CLUSTER"),
					resource.TestCheckResourceAttr(resourceName, "scopes.1.name", userScopeName),
					resource.TestCheckResourceAttr(resourceName, "scopes.1.type", "DATA_LAKE"),
				),
			},
			{
				Config: acc.ConfigDatabaseUserWithScopes(projectID, username, password, "atlasAdmin", nil),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "scopes.#", "0"),
				),
			},
		},
	})
}

func TestAccConfigRSDatabaseUser_updateToEmptyLabels(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		username  = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigDatabaseUserWithLabels(projectID, username, "atlasAdmin", "",
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
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "labels.0.key", "key 1"),
					resource.TestCheckResourceAttr(resourceName, "labels.0.value", "value 1"),
					resource.TestCheckResourceAttr(resourceName, "labels.1.key", "key 2"),
					resource.TestCheckResourceAttr(resourceName, "labels.1.value", "value 2"),
				),
			},
			{
				Config: acc.ConfigDatabaseUserWithLabels(projectID, username, "atlasAdmin", "", nil),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "0"),
				),
			},
		},
	})
}

func TestAccConfigRSDatabaseUser_withLDAPAuthType(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		username  = acc.RandomLDAPName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigDatabaseUserWithLDAPAuthType(projectID, username, "atlasAdmin", "First Key", "First value"),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "ldap_auth_type", "USER"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "$external"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCOnfigRSDatabaseUser_withOIDCAuthType(t *testing.T) {
	var (
		projectID         = acc.ProjectIDExecution(t)
		idpID             = os.Getenv("MONGODB_ATLAS_FEDERATED_IDP_ID")
		workforceAuthType = "IDP_GROUP"
		workloadAuthType  = "USER"
		usernameWorkforce = fmt.Sprintf("%s/%s", idpID, workforceAuthType)
		usernameWorkload  = fmt.Sprintf("%s/%s", idpID, workloadAuthType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigDataBaseUserWithOIDCAuthType(projectID, usernameWorkforce, workforceAuthType, "admin", "atlasAdmin"),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "username", usernameWorkforce),
					resource.TestCheckResourceAttr(resourceName, "oidc_auth_type", workforceAuthType),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
				),
			},
			{
				Config: acc.ConfigDataBaseUserWithOIDCAuthType(projectID, usernameWorkload, workloadAuthType, "$external", "atlasAdmin"),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "username", usernameWorkload),
					resource.TestCheckResourceAttr(resourceName, "oidc_auth_type", workloadAuthType),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "$external"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.Attributes["project_id"] == "" {
			return fmt.Errorf("no project_id is set")
		}
		if rs.Primary.Attributes["auth_database_name"] == "" {
			return fmt.Errorf("no auth_database_name is set")
		}
		if rs.Primary.Attributes["username"] == "" {
			return fmt.Errorf("no username is set")
		}

		authDB := rs.Primary.Attributes["auth_database_name"]
		projectID := rs.Primary.Attributes["project_id"]
		username := rs.Primary.Attributes["username"]

		if _, _, err := acc.ConnV2().DatabaseUsersApi.GetDatabaseUser(context.Background(), projectID, authDB, username).Execute(); err == nil {
			return nil
		}

		return fmt.Errorf("database user(%s-%s-%s) does not exist", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["username"], rs.Primary.Attributes["auth_database_name"])
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_database_user" {
			continue
		}

		projectID, username, authDatabaseName, err := databaseuser.SplitDatabaseUserImportID(rs.Primary.ID)
		if err != nil {
			continue
		}
		// Try to find the database user
		_, _, err = acc.ConnV2().DatabaseUsersApi.GetDatabaseUser(context.Background(), projectID, authDatabaseName, username).Execute()
		if err == nil {
			return fmt.Errorf("database user (%s) still exists", projectID)
		}
	}

	return nil
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		return fmt.Sprintf("%s-%s-%s", ids["project_id"], ids["username"], ids["auth_database_name"]), nil
	}
}

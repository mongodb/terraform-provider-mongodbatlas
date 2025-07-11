package databaseuser_test

import (
	"context"
	"fmt"
	"os"
	"slices"
	"testing"

	"maps"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/databaseuser"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"go.mongodb.org/atlas-sdk/v20250312005/admin"
)

const (
	resourceName         = "mongodbatlas_database_user.test"
	dataSourceName       = "data.mongodbatlas_database_user.test"
	dataSourcePluralName = "data.mongodbatlas_database_users.test"
	dataSourceSingular   = `
		data "mongodbatlas_database_user" "test" {
			username           = mongodbatlas_database_user.test.username
			project_id         = mongodbatlas_database_user.test.project_id
			auth_database_name = mongodbatlas_database_user.test.auth_database_name
		}`
)

var (
	importStep = resource.TestStep{
		ResourceName:            resourceName,
		ImportStateIdFunc:       importStateIDFunc(resourceName),
		ImportState:             true,
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: []string{"password"},
	}
)

func TestAccDatabaseUser_basic(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		username    = acc.RandomName()
		extraChecks = []resource.TestCheckFunc{
			resource.TestCheckNoResourceAttr(resourceName, "description"),
			resource.TestCheckNoResourceAttr(dataSourceName, "description"),
			resource.TestCheckResourceAttr(resourceName, "password", "test-acc-password"),
		}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigDatabaseUserBasic(projectID, username, "atlasAdmin", "First Key", "First value") + dataSourceSingular,
				Check: checkAttrs(
					projectID,
					username,
					"admin",
					map[string]string{
						"labels.#":          "1",
						"labels.0.key":      "First Key",
						"labels.0.value":    "First value",
						"roles.#":           "1",
						"roles.0.role_name": "atlasAdmin",
						"x509_type":         "NONE",
					},
					extraChecks...,
				),
			},
			{
				Config: acc.ConfigDatabaseUserBasic(projectID, username, "read", "Second Key", "Second value") + dataSourceSingular,
				Check: checkAttrs(
					projectID,
					username,
					"admin",
					map[string]string{
						"labels.#":          "1",
						"labels.0.key":      "Second Key",
						"labels.0.value":    "Second value",
						"roles.#":           "1",
						"roles.0.role_name": "read",
						"x509_type":         "NONE",
					},
					extraChecks...,
				),
			},
			importStep,
		},
	})
}

func TestAccDatabaseUser_withX509TypeCustomer(t *testing.T) {
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
				Config: acc.ConfigDatabaseUserWithX509Type(projectID, username, x509Type, "atlasAdmin", "First Key", "First value") + dataSourceSingular,
				Check: checkAttrs(
					projectID,
					username,
					"$external",
					map[string]string{
						"labels.#":       "1",
						"labels.0.key":   "First Key",
						"labels.0.value": "First value",
						"x509_type":      x509Type,
					},
				),
			},
			importStep,
		},
	})
}

func TestAccDatabaseUser_withX509TypeManaged(t *testing.T) {
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
				Config: acc.ConfigDatabaseUserWithX509Type(projectID, username, x509Type, "atlasAdmin", "First Key", "First value") + dataSourceSingular,
				Check: checkAttrs(
					projectID,
					username,
					"$external",
					map[string]string{
						"labels.#":  "1",
						"x509_type": x509Type,
					},
				),
			},
			importStep,
		},
	})
}

func TestAccDatabaseUser_withAWSIAMType(t *testing.T) {
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
				Config: acc.ConfigDatabaseUserWithAWSIAMType(projectID, username, "atlasAdmin", "First Key", "First value") + dataSourceSingular,
				Check: checkAttrs(
					projectID,
					username,
					"$external",
					map[string]string{
						"aws_iam_type": "USER",
					},
				),
			},
			importStep,
		},
	})
}

func TestAccDatabaseUser_withLabelsAndDescription(t *testing.T) {
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
				Config: dataSourceSingular + acc.ConfigDatabaseUserWithLabels(projectID, username, "atlasAdmin", description1, nil),
				Check: checkAttrs(
					projectID,
					username,
					"admin",
					map[string]string{
						"description": description1,
						"labels.#":    "0",
					},
				),
			},
			{
				Config: dataSourceSingular + acc.ConfigDatabaseUserWithLabels(projectID, username, "atlasAdmin", description2,
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
				Check: checkAttrs(
					projectID,
					username,
					"admin",
					map[string]string{
						"description":    description2,
						"labels.#":       "2",
						"labels.0.key":   "key 1",
						"labels.0.value": "value 1",
						"labels.1.key":   "key 2",
						"labels.1.value": "value 2",
					},
				),
			},
			{
				Config: dataSourceSingular + acc.ConfigDatabaseUserWithLabels(projectID, username, "read", "",
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
				Check: checkAttrs(
					projectID,
					username,
					"admin",
					map[string]string{"labels.#": "3"}, // Labels have different order in resource and data source
					resource.TestCheckNoResourceAttr(resourceName, "description"),
				),
			},
			{
				Config: dataSourceSingular + acc.ConfigDatabaseUserWithLabels(projectID, username, "read", "", nil),
				Check: checkAttrs(
					projectID,
					username,
					"admin",
					map[string]string{
						"labels.#": "0",
					},
					resource.TestCheckNoResourceAttr(resourceName, "description"),
				),
			},
			importStep,
		},
	})
}

func TestAccDatabaseUser_withRoles(t *testing.T) {
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
				Config: dataSourceSingular + acc.ConfigDatabaseUserWithRoles(projectID, username, password,
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
				Check: checkAttrs(
					projectID,
					username,
					"admin",
					map[string]string{
						"roles.#":                 "2",
						"roles.0.role_name":       "read",
						"roles.0.database_name":   "admin",
						"roles.0.collection_name": "stir",
						"roles.1.role_name":       "read",
						"roles.1.database_name":   "admin",
						"roles.1.collection_name": "unpledged",
					},
				),
			},
			{
				Config: dataSourceSingular + acc.ConfigDatabaseUserWithRoles(projectID, username, password,
					[]*admin.DatabaseUserRole{
						{
							RoleName:     "read",
							DatabaseName: "admin",
						},
					},
				),
				Check: checkAttrs(
					projectID,
					username,
					"admin",
					map[string]string{
						"roles.#":               "1",
						"roles.0.role_name":     "read",
						"roles.0.database_name": "admin",
					},
				),
			},
			importStep,
		},
	})
}

func TestAccDatabaseUser_withScopes(t *testing.T) {
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
				Config: dataSourceSingular + acc.ConfigDatabaseUserWithScopes(projectID, username, password, "atlasAdmin",
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
				Check: checkAttrs(
					projectID,
					username,
					"admin",
					map[string]string{
						"scopes.#":      "2",
						"scopes.0.name": userScopeName,
						"scopes.0.type": "CLUSTER",
						"scopes.1.name": userScopeName,
						"scopes.1.type": "DATA_LAKE",
					},
				),
			},
			{
				Config: dataSourceSingular + acc.ConfigDatabaseUserWithScopes(projectID, username, password, "atlasAdmin",
					[]*admin.UserScope{
						{
							Name: userScopeName,
							Type: "CLUSTER",
						},
					},
				),
				Check: checkAttrs(
					projectID,
					username,
					"admin",
					map[string]string{
						"scopes.#":      "1",
						"scopes.0.name": userScopeName,
						"scopes.0.type": "CLUSTER",
					},
				),
			},
			{
				Config: dataSourceSingular + acc.ConfigDatabaseUserWithScopes(projectID, username, password, "atlasAdmin", nil),
				Check: checkAttrs(
					projectID,
					username,
					"admin",
					map[string]string{
						"scopes.#": "0",
					},
				),
			},
			importStep,
		},
	})
}

func TestAccDatabaseUser_withLDAPAuthType(t *testing.T) {
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
				Config: acc.ConfigDatabaseUserWithLDAPAuthType(projectID, username, "atlasAdmin", "First Key", "First value") + dataSourceSingular,
				Check: checkAttrs(
					projectID,
					username,
					"$external",
					map[string]string{
						"ldap_auth_type": "USER",
					},
				),
			},
			importStep,
		},
	})
}

func TestAccDatabaseUser_withOIDCAuthType(t *testing.T) {
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
				Config: acc.ConfigDataBaseUserWithOIDCAuthType(projectID, usernameWorkforce, workforceAuthType, "admin", "atlasAdmin") + dataSourceSingular,
				Check: checkAttrs(
					projectID,
					usernameWorkforce,
					"admin",
					map[string]string{
						"oidc_auth_type": workforceAuthType,
					},
				),
			},
			{
				Config: acc.ConfigDataBaseUserWithOIDCAuthType(projectID, usernameWorkload, workloadAuthType, "$external", "atlasAdmin") + dataSourceSingular,
				Check: checkAttrs(
					projectID,
					usernameWorkload,
					"$external",
					map[string]string{
						"oidc_auth_type": workloadAuthType,
					},
				),
			},
			importStep,
		},
	})
}

func checkAttrs(projectID, username, authDBName string, extraAttrs map[string]string, extra ...resource.TestCheckFunc) resource.TestCheckFunc {
	attrsMap := map[string]string{
		"project_id":         projectID,
		"username":           username,
		"auth_database_name": authDBName,
	}
	maps.Copy(attrsMap, extraAttrs)
	check := acc.CheckRSAndDS(resourceName, conversion.Pointer(dataSourceName), nil, nil, attrsMap, extra...)
	checks := slices.Concat(extra, []resource.TestCheckFunc{check, checkExists(resourceName)})
	return resource.ComposeAggregateTestCheckFunc(checks...)
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

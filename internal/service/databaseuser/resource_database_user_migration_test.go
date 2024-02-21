package databaseuser_test

import (
	"os"
	"testing"

	"go.mongodb.org/atlas-sdk/v20231115007/admin"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestAccMigrationConfigRSDatabaseUser_Basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_database_user.basic_ds"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		username     = acc.RandomName()
		config       = acc.ConfigDatabaseUserBasic(projectName, orgID, "atlasAdmin", username, "First Key", "First value")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { mig.PreCheckBasic(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", "test-acc-password"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "1"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestAccMigrationConfigRSDatabaseUser_withX509TypeCustomer(t *testing.T) {
	var (
		resourceName = "mongodbatlas_database_user.test"
		username     = "CN=ellen@example.com,OU=users,DC=example,DC=com"
		x509Type     = "CUSTOMER"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		config       = acc.ConfigDatabaseUserWithX509Type(projectName, orgID, "atlasAdmin", username, "First Key", "First value", x509Type)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { mig.PreCheckBasic(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "x509_type", x509Type),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "$external"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}
func TestAccMigrationConfigRSDatabaseUser_withAWSIAMType(t *testing.T) {
	var (
		resourceName = "mongodbatlas_database_user.test"
		username     = "arn:aws:iam::358363220050:user/mongodb-aws-iam-auth-test-user"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		config       = acc.ConfigDatabaseUserWithAWSIAMType(projectName, orgID, "atlasAdmin", username, "First Key", "First value")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { mig.PreCheckBasic(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "aws_iam_type", "USER"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "$external"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestAccMigrationConfigRSDatabaseUser_withLabels(t *testing.T) {
	var (
		dbUser       admin.CloudDatabaseUser
		resourceName = "mongodbatlas_database_user.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		username     = acc.RandomName()
		config       = acc.ConfigDatabaseUserWithLabels(projectName, orgID, "atlasAdmin", username,
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
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyDatabaseUser,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
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
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestAccMigrationConfigRSDatabaseUser_withEmptyLabels(t *testing.T) {
	var (
		resourceName = "mongodbatlas_database_user.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		username     = acc.RandomName()
		config       = acc.ConfigDatabaseUserWithLabels(projectName, orgID, "atlasAdmin", username, nil)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { mig.PreCheckBasic(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", "test-acc-password"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "0"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestAccMigrationConfigRSDatabaseUser_withRoles(t *testing.T) {
	var (
		resourceName = "mongodbatlas_database_user.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		username     = acc.RandomName()
		password     = acc.RandomName()
		config       = acc.ConfigDatabaseUserWithRoles(username, password, projectName, orgID,
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
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { mig.PreCheckBasic(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "2"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestAccMigrationConfigRSDatabaseUser_withScopes(t *testing.T) {
	var (
		resourceName = "mongodbatlas_database_user.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		clusterName  = acc.RandomClusterName()
		username     = acc.RandomName()
		password     = acc.RandomName()
		config       = acc.ConfigDatabaseUserWithScopes(username, password, "atlasAdmin", projectName, orgID,
			[]*admin.UserScope{
				{
					Name: clusterName,
					Type: "CLUSTER",
				},
			},
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { mig.PreCheckBasic(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "scopes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scopes.0.name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "scopes.0.type", "CLUSTER"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestAccMigrationConfigRSDatabaseUser_withEmptyScopes(t *testing.T) {
	var (
		resourceName = "mongodbatlas_database_user.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		username     = acc.RandomName()
		password     = acc.RandomName()
		config       = acc.ConfigDatabaseUserWithScopes(username, password, "atlasAdmin", projectName, orgID, nil)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { mig.PreCheckBasic(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "scopes.#", "0"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestAccMigrationConfigRSDatabaseUser_withLDAPAuthType(t *testing.T) {
	var (
		resourceName = "mongodbatlas_database_user.test"
		username     = "CN=david@example.com,OU=users,DC=example,DC=com"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		config       = acc.ConfigDatabaseUserWithLDAPAuthType(projectName, orgID, "atlasAdmin", username, "First Key", "First value")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { mig.PreCheckBasic(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "ldap_auth_type", "USER"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "$external"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

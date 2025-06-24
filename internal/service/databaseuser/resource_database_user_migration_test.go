package databaseuser_test

import (
	"testing"

	"go.mongodb.org/atlas-sdk/v20250312004/admin"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigConfigRSDatabaseUser_Basic(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		username  = acc.RandomName()
		config    = acc.ConfigDatabaseUserBasic(projectID, username, "atlasAdmin", "First Key", "First value")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { mig.PreCheckBasic(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
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

func TestMigConfigRSDatabaseUser_withX509TypeCustomer(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		username  = acc.RandomLDAPName()
		x509Type  = "CUSTOMER"
		config    = acc.ConfigDatabaseUserWithX509Type(projectID, username, x509Type, "atlasAdmin", "First Key", "First value")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { mig.PreCheckBasic(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeAggregateTestCheckFunc(
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
func TestMigConfigRSDatabaseUser_withAWSIAMType(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		username  = acc.RandomIAMUser()
		config    = acc.ConfigDatabaseUserWithAWSIAMType(projectID, username, "atlasAdmin", "First Key", "First value")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { mig.PreCheckBasic(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
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

func TestMigConfigRSDatabaseUser_withLabels(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		username  = acc.RandomName()
		config    = acc.ConfigDatabaseUserWithLabels(projectID, username, "atlasAdmin", "",
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
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
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

func TestMigConfigRSDatabaseUser_withEmptyLabels(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		username  = acc.RandomName()
		config    = acc.ConfigDatabaseUserWithLabels(projectID, username, "atlasAdmin", "", nil)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { mig.PreCheckBasic(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
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

func TestMigConfigRSDatabaseUser_withRoles(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		username  = acc.RandomName()
		password  = acc.RandomName()
		config    = acc.ConfigDatabaseUserWithRoles(projectID, username, password,
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
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
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

func TestMigConfigRSDatabaseUser_withScopes(t *testing.T) {
	var (
		projectID     = acc.ProjectIDExecution(t)
		userScopeName = acc.RandomName()
		username      = acc.RandomName()
		password      = acc.RandomName()
		config        = acc.ConfigDatabaseUserWithScopes(projectID, username, password, "atlasAdmin",
			[]*admin.UserScope{
				{
					Name: userScopeName,
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
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "scopes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scopes.0.name", userScopeName),
					resource.TestCheckResourceAttr(resourceName, "scopes.0.type", "CLUSTER"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestMigConfigRSDatabaseUser_withEmptyScopes(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		username  = acc.RandomName()
		password  = acc.RandomName()
		config    = acc.ConfigDatabaseUserWithScopes(projectID, username, password, "atlasAdmin", nil)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { mig.PreCheckBasic(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeAggregateTestCheckFunc(
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

func TestMigConfigRSDatabaseUser_withLDAPAuthType(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		username  = acc.RandomLDAPName()
		config    = acc.ConfigDatabaseUserWithLDAPAuthType(projectID, username, "atlasAdmin", "First Key", "First value")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { mig.PreCheckBasic(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
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

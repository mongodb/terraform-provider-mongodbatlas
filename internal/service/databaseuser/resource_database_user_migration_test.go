package databaseuser_test

import (
	"os"
	"testing"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestAccMigrationConfigRSDatabaseUser_Basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_database_user.basic_ds"
		username     = acctest.RandomWithPrefix("dbUser")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { mig.PreCheckBasic(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(t),
				Config:            acc.ConfigDatabaseUserBasic(projectName, orgID, "atlasAdmin", username, "First Key", "First value"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", "test-acc-password"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "1"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   acc.ConfigDatabaseUserBasic(projectName, orgID, "atlasAdmin", username, "First Key", "First value"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}

func TestAccMigrationConfigRSDatabaseUser_WithX509TypeCustomer(t *testing.T) {
	var (
		resourceName = "mongodbatlas_database_user.test"
		username     = "CN=ellen@example.com,OU=users,DC=example,DC=com"
		x509Type     = "CUSTOMER"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { mig.PreCheckBasic(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(t),
				Config:            acc.ConfigDatabaseUserWithX509Type(projectName, orgID, "atlasAdmin", username, "First Key", "First value", x509Type),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "x509_type", x509Type),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "$external"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   acc.ConfigDatabaseUserWithX509Type(projectName, orgID, "atlasAdmin", username, "First Key", "First value", x509Type),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}
func TestAccMigrationConfigRSDatabaseUser_WithAWSIAMType(t *testing.T) {
	var (
		resourceName = "mongodbatlas_database_user.test"
		username     = "arn:aws:iam::358363220050:user/mongodb-aws-iam-auth-test-user"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { mig.PreCheckBasic(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(t),
				Config:            acc.ConfigDatabaseUserWithAWSIAMType(projectName, orgID, "atlasAdmin", username, "First Key", "First value"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "aws_iam_type", "USER"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "$external"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   acc.ConfigDatabaseUserWithAWSIAMType(projectName, orgID, "atlasAdmin", username, "First Key", "First value"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}

func TestAccMigrationConfigRSDatabaseUser_WithLabels(t *testing.T) {
	var (
		dbUser       matlas.DatabaseUser
		resourceName = "mongodbatlas_database_user.test"
		username     = acctest.RandomWithPrefix("test-acc")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyDatabaseUser,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(t),
				Config: acc.ConfigDatabaseUserWithLabels(projectName, orgID, "atlasAdmin", username,
					[]matlas.Label{
						{
							Key:   "key 1",
							Value: "value 1",
						},
						{
							Key:   "key 2",
							Value: "value 2",
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
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config: acc.ConfigDatabaseUserWithLabels(projectName, orgID, "atlasAdmin", username,
					[]matlas.Label{
						{
							Key:   "key 1",
							Value: "value 1",
						},
						{
							Key:   "key 2",
							Value: "value 2",
						},
					},
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}
func TestAccMigrationConfigRSDatabaseUser_WithEmptyLabels(t *testing.T) {
	var (
		resourceName = "mongodbatlas_database_user.test"
		username     = acctest.RandomWithPrefix("test-acc")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { mig.PreCheckBasic(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(t),
				Config:            acc.ConfigDatabaseUserWithLabels(projectName, orgID, "atlasAdmin", username, []matlas.Label{}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", "test-acc-password"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "0"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   acc.ConfigDatabaseUserWithLabels(projectName, orgID, "atlasAdmin", username, []matlas.Label{}),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}

func TestAccMigrationConfigRSDatabaseUser_WithRoles(t *testing.T) {
	var (
		resourceName = "mongodbatlas_database_user.test"
		username     = acctest.RandomWithPrefix("test-acc-user-")
		password     = acctest.RandomWithPrefix("test-acc-pass-")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { mig.PreCheckBasic(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(t),
				Config: acc.ConfigDatabaseUserWithRoles(username, password, projectName, orgID,
					[]*matlas.Role{
						{
							RoleName:       "read",
							DatabaseName:   "admin",
							CollectionName: "stir",
						},
						{
							RoleName:       "read",
							DatabaseName:   "admin",
							CollectionName: "unpledged",
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "2"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config: acc.ConfigDatabaseUserWithRoles(username, password, projectName, orgID,
					[]*matlas.Role{
						{
							RoleName:       "read",
							DatabaseName:   "admin",
							CollectionName: "stir",
						},
						{
							RoleName:       "read",
							DatabaseName:   "admin",
							CollectionName: "unpledged",
						},
					},
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}

func TestAccMigrationConfigRSDatabaseUser_WithScopes(t *testing.T) {
	var (
		resourceName = "mongodbatlas_database_user.test"
		username     = acctest.RandomWithPrefix("test-acc-user-")
		password     = acctest.RandomWithPrefix("test-acc-pass-")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		clusterName  = acctest.RandomWithPrefix("test-acc-cluster")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { mig.PreCheckBasic(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(t),
				Config: acc.ConfigDatabaseUserWithScopes(username, password, projectName, orgID, "atlasAdmin", clusterName,
					[]*matlas.Scope{
						{
							Name: "test-acc-nurk4llu2z",
							Type: "CLUSTER",
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "scopes.#", "1"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config: acc.ConfigDatabaseUserWithScopes(username, password, projectName, orgID, "atlasAdmin", clusterName,
					[]*matlas.Scope{
						{
							Name: "test-acc-nurk4llu2z",
							Type: "CLUSTER",
						},
					},
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}

func TestAccMigrationConfigRSDatabaseUser_WithScopesAndEmpty(t *testing.T) {
	var (
		resourceName = "mongodbatlas_database_user.test"
		username     = acctest.RandomWithPrefix("test-acc-user-")
		password     = acctest.RandomWithPrefix("test-acc-pass-")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		clusterName  = acctest.RandomWithPrefix("test-acc-cluster")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { mig.PreCheckBasic(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(t),
				Config: acc.ConfigDatabaseUserWithScopes(username, password, projectName, orgID, "atlasAdmin", clusterName,
					[]*matlas.Scope{},
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "scopes.#", "0"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config: acc.ConfigDatabaseUserWithScopes(username, password, projectName, orgID, "atlasAdmin", clusterName,
					[]*matlas.Scope{},
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}

func TestAccMigrationConfigRSDatabaseUser_WithLDAPAuthType(t *testing.T) {
	var (
		resourceName = "mongodbatlas_database_user.test"
		username     = "CN=david@example.com,OU=users,DC=example,DC=com"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { mig.PreCheckBasic(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(t),
				Config:            acc.ConfigDatabaseUserWithLDAPAuthType(projectName, orgID, "atlasAdmin", username, "First Key", "First value"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "ldap_auth_type", "USER"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "$external"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   acc.ConfigDatabaseUserWithLDAPAuthType(projectName, orgID, "atlasAdmin", username, "First Key", "First value"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}

package mongodbatlas

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas/testutils"
)

func TestAccFederatedDatabaseInstance_MigraitonTest_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_federated_database_instance.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		name         = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckBasic(t) },
		CheckDestroy: testAccCheckMongoDBAtlasFederatedDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: "1.11.0",
						Source:            "mongodb/mongodbatlas",
					},
				},
				Config: testAccMongoDBAtlasFederatedDatabaseInstanceConfigFirstSteps(name, projectName, orgID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{
				ProviderFactories: testAccProviderFactories,
				Config:            testAccMongoDBAtlasFederatedDatabaseInstanceConfigFirstSteps(name, projectName, orgID),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						testutils.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
			{
				ResourceName:      resourceName,
				ProviderFactories: testAccProviderFactories,
				ImportStateIdFunc: testAccCheckMongoDBAtlasFederatedDatabaseInstanceImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccFederatedDatabaseInstance_Migration_S3bucket(t *testing.T) {
	SkipTestExtCred(t)
	var (
		resourceName = "mongodbatlas_federated_database_instance.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		name         = acctest.RandomWithPrefix("test-acc")
		policyName   = acctest.RandomWithPrefix("test-acc")
		roleName     = acctest.RandomWithPrefix("test-acc")
		testS3Bucket = os.Getenv("AWS_S3_BUCKET")
		region       = "VIRGINIA_USA"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckBasic(t) },
		CheckDestroy: testAccCheckMongoDBAtlasFederatedDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"aws": {
						VersionConstraint: "5.1.0",
						Source:            "hashicorp/aws",
					},
					"mongodbatlas": {
						VersionConstraint: "1.11.0",
						Source:            "mongodb/mongodbatlas",
					},
				},
				Config: testAccMongoDBAtlasFederatedDatabaseInstanceConfigS3Bucket(policyName, roleName, projectName, orgID, name, testS3Bucket, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"aws": {
						VersionConstraint: "5.1.0",
						Source:            "hashicorp/aws",
					},
				},
				ProviderFactories: testAccProviderFactories,
				Config:            testAccMongoDBAtlasFederatedDatabaseInstanceConfigS3Bucket(policyName, roleName, projectName, orgID, name, testS3Bucket, region),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						testutils.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
			{
				ResourceName:      resourceName,
				ProviderFactories: testAccProviderFactories,
				ImportStateIdFunc: testAccCheckMongoDBAtlasFederatedDatabaseInstanceImportStateIDFuncS3Bucket(resourceName, testS3Bucket),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccFederatedDatabaseInstance_Migration_atlasCluster(t *testing.T) {
	var (
		resourceName = "mongodbatlas_federated_database_instance.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		name         = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckBasic(t) },
		CheckDestroy: testAccCheckMongoDBAtlasFederatedDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: "1.11.0",
						Source:            "mongodb/mongodbatlas",
					},
				},
				Config: testAccMongoDBAtlasFederatedDatabaseInstanceAtlasProviderConfig(projectName, orgID, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{
				ProviderFactories: testAccProviderFactories,
				Config:            testAccMongoDBAtlasFederatedDatabaseInstanceAtlasProviderConfig(projectName, orgID, name),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						testutils.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}

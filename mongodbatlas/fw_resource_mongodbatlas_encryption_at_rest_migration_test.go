package mongodbatlas

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/mwielbut/pointy"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccAdvRS_Migration_EncryptionAtRest_basicAWS(t *testing.T) {
	var (
		resourceName = "mongodbatlas_encryption_at_rest.test"
		// projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		projectID = "63bec58b014da65b8f73c06c"

		awsKms = matlas.AwsKms{
			Enabled:             pointy.Bool(true),
			CustomerMasterKeyID: os.Getenv("AWS_CUSTOMER_MASTER_KEY_ID"),
			Region:              os.Getenv("AWS_REGION"),
			RoleID:              os.Getenv("AWS_ROLE_ID"),
		}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testCheckAwsEnv(t) },
		CheckDestroy: testAccCheckMongoDBAtlasEncryptionAtRestDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: "1.11.0",
						Source:            "mongodb/mongodbatlas",
					},
				},
				Config: testAccMongoDBAtlasEncryptionAtRestConfigAwsKms(projectID, &awsKms),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEncryptionAtRestExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				ProtoV6ProviderFactories: testAccProviderV6Factories,
				Config:                   testAccMongoDBAtlasEncryptionAtRestConfigAwsKms(projectID, &awsKms),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}

func TestAccAdvRS_Migration_EncryptionAtRest_WithRole_basicAWS(t *testing.T) {
	// SkipTest(t) // For now it will skipped because of aws errors reasons, already made another test using terratest.
	// SkipTestExtCred(t)
	var (
		resourceName = "mongodbatlas_encryption_at_rest.test"
		// projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		projectID   = "63bec58b014da65b8f73c06c"
		accessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
		secretKey   = os.Getenv("AWS_SECRET_ACCESS_KEY")
		policyName  = acctest.RandomWithPrefix("test-aws-policy")
		roleName    = acctest.RandomWithPrefix("test-aws-role")

		awsKms = matlas.AwsKms{
			Enabled:             pointy.Bool(true),
			CustomerMasterKeyID: os.Getenv("AWS_CUSTOMER_MASTER_KEY_ID"),
			Region:              os.Getenv("AWS_REGION"),
		}
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testCheckAwsEnv(t) },
		CheckDestroy: testAccCheckMongoDBAtlasEncryptionAtRestDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: "1.11.0",
						Source:            "mongodb/mongodbatlas",
					},
					"aws": {
						VersionConstraint: "5.1.0",
						Source:            "hashicorp/aws",
					},
				},
				Config: testAccMongoDBAtlasEncryptionAtRestConfigAwsKmsWithRole(awsKms.Region, accessKeyID, secretKey, projectID, policyName, roleName, false, &awsKms),
			},
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"aws": {
						VersionConstraint: "5.1.0",
						Source:            "hashicorp/aws",
					},
				},
				ProtoV6ProviderFactories: testAccProviderV6Factories,
				Config:                   testAccMongoDBAtlasEncryptionAtRestConfigAwsKmsWithRole(awsKms.Region, accessKeyID, secretKey, projectID, policyName, roleName, false, &awsKms),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEncryptionAtRestExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}

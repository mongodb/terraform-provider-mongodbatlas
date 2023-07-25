package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccConfigRSCloudProviderAccessSetupAWS_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_cloud_provider_access_setup.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		targetRole   = matlas.CloudProviderAccessRole{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		// same as regular cloud provider access resource
		CheckDestroy: testAccCheckMongoDBAtlasProviderAccessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudProviderAccessSetupAWS(orgID, projectName),
				Check: resource.ComposeTestCheckFunc(
					// same as regular cloud resource
					testAccCheckMongoDBAtlasProviderAccessExists(resourceName, &targetRole),
					resource.TestCheckResourceAttrSet(resourceName, "aws.atlas_assumed_role_external_id"),
					resource.TestCheckResourceAttrSet(resourceName, "aws.atlas_aws_account_arn"),
				),
			},
		},
	},
	)
}

func TestAccConfigRSCloudProviderAccessSetupAWS_importBasic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_cloud_provider_access_setup.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		targetRole   = matlas.CloudProviderAccessRole{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasProviderAccessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudProviderAccessSetupAWS(orgID, projectName),
				Check: resource.ComposeTestCheckFunc(
					// same as regular cloud provider because we are just checking in the api
					testAccCheckMongoDBAtlasProviderAccessExists(resourceName, &targetRole),
					resource.TestCheckResourceAttrSet(resourceName, "aws.atlas_assumed_role_external_id"),
					resource.TestCheckResourceAttrSet(resourceName, "aws.atlas_aws_account_arn"),
				),
			},
			{
				ResourceName: resourceName,
				// ID remains the same project-id, provider-name and id for consistency
				ImportStateIdFunc: testAccCheckMongoDBAtlasCloudProviderAccessImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	},
	)
}

// atlas_azure_app_id = "6f2deb0d-be72-4524-a403-df531868bac0"
// service_principal_id = "48f1d2a6-d0e9-482a-83a4-b8dd7dddc5c1"
// tenant_id = "91405384-d71e-47f5-92dd-759e272cdc1c"
func TestAccConfigRSCloudProviderAccessSetupAzure_basic(t *testing.T) {
	var (
		resourceName       = "mongodbatlas_cloud_provider_access_setup.test"
		orgID              = os.Getenv("MONGODB_ATLAS_ORG_ID")
		atlasAzureAppID    = "6f2deb0d-be72-4524-a403-df531868bac0" // os.Getenv("AZURE_ATLAS_APP_ID")
		servicePrincipalID = "48f1d2a6-d0e9-482a-83a4-b8dd7dddc5c1" // os.Getenv("AZURE_SERVICE_PRICIPAL_ID")
		tenantID           = "91405384-d71e-47f5-92dd-759e272cdc1c" // os.Getenv("AZURE_TENANT_OD")
		projectName        = acctest.RandomWithPrefix("test-acc")
		targetRole         = matlas.CloudProviderAccessRole{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudProviderAccessSetupAzure(orgID, projectName, atlasAzureAppID, servicePrincipalID, tenantID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProviderAccessExists(resourceName, &targetRole),
					resource.TestCheckResourceAttrSet(resourceName, "role_id"),
				),
			},
		},
	},
	)
}

func testAccMongoDBAtlasCloudProviderAccessSetupAWS(orgID, projectName string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_project" "test" {
		name   = %[2]q
		org_id = %[1]q
	}
	resource "mongodbatlas_cloud_provider_access_setup" "test" {
		project_id = mongodbatlas_project.test.id
		provider_name = "AWS"
	 }
	`, orgID, projectName)
}

func testAccMongoDBAtlasCloudProviderAccessSetupAzure(orgID, projectName, atlasAzureAppID, servicePrincipalID, tenantID string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_project" "test" {
		name   = %[2]q
		org_id = %[1]q
	}
	resource "mongodbatlas_cloud_provider_access_setup" "test" {
		project_id = mongodbatlas_project.test.id
		provider_name = "AZURE"
		azure_config {
			atlas_azure_app_id = %[3]q
			service_principal_id = %[4]q
			tenant_id = %[5]q
		}
	 }
	`, orgID, projectName, atlasAzureAppID, servicePrincipalID, tenantID)
}

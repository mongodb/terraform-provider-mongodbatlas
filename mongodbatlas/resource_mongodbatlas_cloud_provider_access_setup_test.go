package mongodbatlas_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccConfigRSCloudProviderAccessSetupAWS_basic(t *testing.T) {
	var (
		resourceName   = "mongodbatlas_cloud_provider_access_setup.test"
		dataSourceName = "data.mongodbatlas_cloud_provider_access_setup.test"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acctest.RandomWithPrefix("test-acc")
		targetRole     = matlas.CloudProviderAccessRole{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasProviderAccessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudProviderAccessSetupAWS(orgID, projectName),
				Check: resource.ComposeTestCheckFunc(
					// same as regular cloud resource
					testAccCheckMongoDBAtlasProviderAccessExists(resourceName, &targetRole),
					resource.TestCheckResourceAttrSet(dataSourceName, "aws_config.0.atlas_assumed_role_external_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "aws_config.0.atlas_aws_account_arn"),
					resource.TestCheckResourceAttrSet(dataSourceName, "created_date"),
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
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasProviderAccessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudProviderAccessSetupAWS(orgID, projectName),
				Check: resource.ComposeTestCheckFunc(
					// same as regular cloud provider because we are just checking in the api
					testAccCheckMongoDBAtlasProviderAccessExists(resourceName, &targetRole),
					resource.TestCheckResourceAttrSet(resourceName, "aws_config.0.atlas_assumed_role_external_id"),
					resource.TestCheckResourceAttrSet(resourceName, "aws_config.0.atlas_aws_account_arn"),
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

func TestAccConfigRSCloudProviderAccessSetupAzure_basic(t *testing.T) {
	var (
		resourceName       = "mongodbatlas_cloud_provider_access_setup.test"
		dataSourceName     = "data.mongodbatlas_cloud_provider_access_setup.test"
		orgID              = os.Getenv("MONGODB_ATLAS_ORG_ID")
		atlasAzureAppID    = os.Getenv("AZURE_ATLAS_APP_ID")
		servicePrincipalID = os.Getenv("AZURE_SERVICE_PRINCIPAL_ID")
		tenantID           = os.Getenv("AZURE_TENANT_ID")
		projectName        = acctest.RandomWithPrefix("test-acc")
		targetRole         = matlas.CloudProviderAccessRole{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckCloudProviderAccessAzure(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudProviderAccessSetupAzure(orgID, projectName, atlasAzureAppID, servicePrincipalID, tenantID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProviderAccessExists(resourceName, &targetRole),
					resource.TestCheckResourceAttrSet(resourceName, "role_id"),
					resource.TestCheckResourceAttrSet(resourceName, "azure_config.0.atlas_azure_app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "azure_config.0.service_principal_id"),
					resource.TestCheckResourceAttrSet(resourceName, "azure_config.0.tenant_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_date"),
					resource.TestCheckResourceAttrSet(resourceName, "last_updated_date"),
					resource.TestCheckResourceAttrSet(dataSourceName, "role_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "azure_config.0.atlas_azure_app_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "azure_config.0.service_principal_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "azure_config.0.tenant_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "created_date"),
					resource.TestCheckResourceAttrSet(dataSourceName, "last_updated_date"),
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

	 data "mongodbatlas_cloud_provider_access_setup" "test" {
		project_id = mongodbatlas_cloud_provider_access_setup.test.project_id
		provider_name = "AWS"
		role_id =  mongodbatlas_cloud_provider_access_setup.test.role_id
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

	 data "mongodbatlas_cloud_provider_access_setup" "test" {
		project_id = mongodbatlas_cloud_provider_access_setup.test.project_id
		provider_name = "AWS"
		role_id =  mongodbatlas_cloud_provider_access_setup.test.role_id
	 }
	`, orgID, projectName, atlasAzureAppID, servicePrincipalID, tenantID)
}
